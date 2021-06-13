package lvm

import (
	"encoding/json"
	"fmt"
	"github.com/google/martian/log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// LogicalVolume specifies attributes of a given lv exists on node.
type LogicalVolume struct {

	// Name of the lvm logical volume
	Name string `json:"name"`

	// Full name of the lvm logical volume
	FullName string `json:"fullName"`

	// UUID denotes a unique identity of a lvm logical volume.
	UUID string `json:"uuid"`

	// Size specifies the total size of logical volume
	Size int64 `json:"size"`

	// Path specifies LVM logical volume path
	Path string `json:"path"`

	// DMPath specifies device mapper path
	DMPath string `json:"dmPath"`

	// LVM logical volume device
	Device string `json:"device"`

	// Name of the VG in which LVM logical volume is created
	VGName string `json:"vgName"`

	//Active state of the LV.
	Active string
}

// VolumeGroup specifies attributes of a given vg exists on node.
type VolumeGroup struct {

	// Name of the lvm volume group
	Name string

	// UUID of the lvm vg
	UUID string

	// Allocation policy of lvm vg
	AllocationPolicy string

	// Format of the lvm vg
	Format string

	// Free size of lvm vg in bytes
	FreeSize int64

	// Total size of lvm vg in bytes
	TotalSize int64

	// Attributes of lvm vg
	Attr string

	// Permission of lvm vg
	Permission string

	// Number of PVs in VG.
	PVCount int32

	// Count of lv created on vg
	LVCount int32

	// Size of Physical Extents in current units.
	ExtendSize int32

	//Total number of Physical Extents.
	ExtendCount int32

	//Total number of unallocated Physical Extents.
	FreeCount int32

	//Maximum number of LVs allowed in VG or 0 if unlimited.
	MaxLV int32

	//Maximum number of PVs allowed in VG or 0 if unlimited.
	MaxPV int32

	//Number of PVs in VG which are missing.
	MissingPVCount int32

	// Number of snapshots
	SnapCount int32
}

// VolumeGroup specifies attributes of a given vg exists on node.
type PhysicalVolume struct {

	// Name of the lvm physical volume
	Name string

	// UUID of the lvm pv
	UUID string

	// Allocation policy of lvm pv
	AllocationPolicy string

	// Format of the lvm pv
	Format string

	// Total amount of unallocated space in current units.
	FreeSize int64

	// Size of PV in current units.
	TotalSize int64

	//Total amount of allocated space in current units.
	UsedSize int64

	// Size of underlying device in current units.
	DevSize int64

	// Attributes of lvm pv
	Attr string

	//Set if PV is used.
	PVInUse string
}

type LVMInterface interface {
	ReloadLVMMetadataCache() error
	getSymLinks(path string) (string, error)
	execLVMCommand(cmnd string, args ...string) ([]byte, error)
	ListLVMVolumeGroup() ([]VolumeGroup, error)
	ListLVMLogicalVolume() ([]LogicalVolume, error)
	ListLVMPhysicalVolume() ([]PhysicalVolume, error)
}

type LVMstruct struct {
	LVMInterface
}

var (
	LVS    = "lvs"
	VGS    = "vgs"
	PVS    = "pvs"
	PVSCAN = "pvscan"

	LVMClient LVMInterface = &LVMstruct{}
)

/*
Function to get LVM Logical volume device
*/
func getLvDeviceName(path string) (string, error) {
	symLink, err := LVMClient.getSymLinks(path)
	if err != nil {
		log.Errorf("lvm: error in getting device name")
		return "", err
	}
	deviceName := strings.Split(symLink, "/")
	return deviceName[len(deviceName)-1], nil
}

func (lvmC *LVMstruct) getSymLinks(path string) (string, error) {
	return filepath.EvalSymlinks(path)
}

// ReloadLVMMetadataCache refreshes lvmetad daemon cache used for
// serving vgs or other lvm utility.
func (lvmC *LVMstruct) ReloadLVMMetadataCache() error {
	args := []string{"--cache"}
	output, err := LVMClient.execLVMCommand(PVSCAN, args...)
	if err != nil {
		log.Errorf("lvm: reload lvm metadata cache: %v - %v", string(output), err)
		return err
	}
	return nil
}

/*
To parse the output of pvs command and store it in PhysicalVolume
*/
func parsePhysicalVolume(m map[string]string) (PhysicalVolume, error) {
	var pv PhysicalVolume
	var err error
	pv.Name = m["pv_name"]
	pv.UUID = m["pv_uuid"]
	pv.Attr = m["pv_attr"]
	pv.PVInUse = m["pv_in_use"]

	int64Map := map[string]*int64{
		"pv_size":  &pv.TotalSize,
		"pv_free":  &pv.FreeSize,
		"pv_used":  &pv.UsedSize,
		"dev_size": &pv.DevSize,
	}

	for k, v := range int64Map {
		sizeBytes, err := strconv.ParseInt(strings.TrimSuffix(strings.ToLower(m[k]), "b"), 10, 64)
		if err != nil {
			log.Errorf("invalid format of %v=%v for lv %v: %v", k, m[k], pv.Name, err)
			return PhysicalVolume{}, err
		}
		*v = sizeBytes
	}

	return pv, err
}

/*
Decode json format and store physical volumes in map[string]string
*/
func decodePvsJSON(raw []byte) ([]PhysicalVolume, error) {
	output := &struct {
		Report []struct {
			PhysicalVolume []map[string]string `json:"pv"`
		} `json:"report"`
	}{}
	var err error
	if err = json.Unmarshal(raw, output); err != nil {
		return nil, err
	}

	if len(output.Report) != 1 {
		return nil, fmt.Errorf("expected exactly one lvm report")
	}

	items := output.Report[0].PhysicalVolume
	pvs := make([]PhysicalVolume, 0, len(items))
	for _, item := range items {
		var pv PhysicalVolume
		if pv, err = parsePhysicalVolume(item); err != nil {
			return pvs, err
		}
		pvs = append(pvs, pv)
	}
	return pvs, nil
}

/*
ListLVMPhysicalVolume invokes `pvs` to list all the available LVM physcial volumes in the node.
*/
func (lvmC *LVMstruct) ListLVMPhysicalVolume() ([]PhysicalVolume, error) {
	if err := LVMClient.ReloadLVMMetadataCache(); err != nil {
		return nil, err
	}

	args := []string{
		"--options", "pv_all",
		"--reportformat", "json",
		"--units", "b",
	}
	output, err := LVMClient.execLVMCommand(PVS, args...)
	if err != nil {
		log.Errorf("lvm: list logical volume cmd %v: %v", args, err)
		return nil, err
	}
	return decodePvsJSON(output)
}

/*
To parse the output of vgs command and store it in VolumeGroup
*/
func parseVolumeGroup(m map[string]string) (VolumeGroup, error) {
	var vg VolumeGroup
	var err error
	vg.Name = m["vg_name"]
	vg.UUID = m["vg_uuid"]
	vg.Attr = m["vg_attr"]

	int64Map := map[string]*int64{
		"vg_size": &vg.TotalSize,
		"vg_free": &vg.FreeSize,
	}

	for k, v := range int64Map {
		sizeBytes, err := strconv.ParseInt(strings.TrimSuffix(strings.ToLower(m[k]), "b"), 10, 64)
		if err != nil {
			log.Errorf("invalid format of %v=%v for lv %v: %v", k, m[k], vg.Name, err)
			return VolumeGroup{}, err
		}
		*v = sizeBytes
	}

	vg.AllocationPolicy = m["vg_allocation_policy"]
	vg.Format = m["vg_fmt"]
	vg.Permission = m["vg_permissions"]

	int32Map := map[string]*int32{
		"lv_count":            &vg.LVCount,
		"pv_count":            &vg.PVCount,
		"max_lv":              &vg.MaxLV,
		"max_pv":              &vg.MaxPV,
		"vg_missing_pv_count": &vg.MissingPVCount,
		"snap_count":          &vg.SnapCount,
	}

	for k, v := range int32Map {
		count, err := strconv.Atoi(m[k])
		if err != nil {
			log.Errorf("invalid format of %v=%v for vg %v: %v", k, m[k], vg.LVCount, err)
			return VolumeGroup{}, err
		}
		*v = int32(count)
	}

	return vg, err
}

/*
Decode json format and store LVM volume group in VolumeGroup
*/
func decodeVgsJSON(raw []byte) ([]VolumeGroup, error) {
	output := &struct {
		Report []struct {
			VolumeGroups []map[string]string `json:"vg"`
		} `json:"report"`
	}{}
	var err error
	if err = json.Unmarshal(raw, output); err != nil {
		return nil, err
	}

	if len(output.Report) != 1 {
		return nil, fmt.Errorf("expected exactly one lvm report")
	}

	items := output.Report[0].VolumeGroups
	vgs := make([]VolumeGroup, 0, len(items))
	for _, item := range items {
		var vg VolumeGroup
		if vg, err = parseVolumeGroup(item); err != nil {
			return vgs, err
		}
		vgs = append(vgs, vg)
	}
	return vgs, nil
}

/*
ListLVMVolumeGroup invokes `vgs` to list all the available volume
groups in the node.
*/
func (lvmC *LVMstruct) ListLVMVolumeGroup() ([]VolumeGroup, error) {
	if err := LVMClient.ReloadLVMMetadataCache(); err != nil {
		return nil, err
	}

	args := []string{
		"--options", "vg_all",
		"--reportformat", "json",
		"--units", "b",
	}
	output, err := LVMClient.execLVMCommand(VGS, args...)
	if err != nil {
		log.Errorf("lvm: list volume group cmd %v: %v", args, err)
		return nil, err
	}
	return decodeVgsJSON(output)
}

/*
To parse the output of lvs command and store it in LogicalVolume
*/
func parseLogicalVolume(m map[string]string) (LogicalVolume, error) {
	var lv LogicalVolume
	var err error

	lv.Name = m["lv_name"]
	lv.FullName = m["lv_full_name"]
	lv.UUID = m["lv_uuid"]
	lv.Path = m["lv_path"]
	lv.DMPath = m["lv_dm_path"]
	lv.VGName = m["vg_name"]
	sizeBytes, err := strconv.ParseInt(strings.TrimSuffix(strings.ToLower(m["lv_size"]), "b"), 10, 64)

	if err != nil {
		err = fmt.Errorf("invalid format of lv_size=%v for lv %v: %v", m["lv_size"], lv.Name, err)
	}
	lv.Size = sizeBytes
	lv.Active = m["lv_active"]
	return lv, err
}

/*
Decode json format and store logical volumes in map[string]string
*/
func decodeLvsJSON(raw []byte) ([]LogicalVolume, error) {
	output := &struct {
		Report []struct {
			LogicalVolumes []map[string]string `json:"lv"`
		} `json:"report"`
	}{}
	var err error
	if err = json.Unmarshal(raw, output); err != nil {
		return nil, err
	}

	if len(output.Report) != 1 {
		return nil, fmt.Errorf("expected exactly one lvm report")
	}

	items := output.Report[0].LogicalVolumes
	lvs := make([]LogicalVolume, 0, len(items))
	for _, item := range items {
		var lv LogicalVolume
		if lv, err = parseLogicalVolume(item); err != nil {
			return lvs, err
		}
		deviceName, err := getLvDeviceName(lv.Path)
		if err != nil {
			log.Errorf("Error: %v", err)
			return nil, err
		}
		lv.Device = deviceName
		lvs = append(lvs, lv)
	}
	return lvs, nil
}

/*
ListLVMLogicalVolume invokes `lvs` to list all the available logical volumes in the node.
*/
func (lvmC *LVMstruct) ListLVMLogicalVolume() ([]LogicalVolume, error) {
	if err := LVMClient.ReloadLVMMetadataCache(); err != nil {
		return nil, err
	}

	args := []string{
		"--options", "lv_all,vg_name",
		"--reportformat", "json",
		"--units", "b",
	}
	output, err := LVMClient.execLVMCommand(LVS, args...)
	if err != nil {
		log.Errorf("lvm: list logical volume cmd %v: %v", args, err)
		return nil, err
	}
	return decodeLvsJSON(output)
}

/*
Wrapper over exec.Command(...)
*/
func (lvmC *LVMstruct) execLVMCommand(cmnd string, args ...string) ([]byte, error) {
	cmd := exec.Command(cmnd, args...)
	return cmd.CombinedOutput()
}
