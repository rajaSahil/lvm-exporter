package lvm

import (
	"errors"
	"reflect"
	"testing"
)

type MockLVMStruct struct {
}

var (
	GetexecLVMCommandFunc         func(string, ...string) ([]byte, error)
	GetReloadLVMMetadataCacheFunc func() error
	GetListLVMVolumeGroupFunc     func() ([]VolumeGroup, error)
	GetListLVMLogicalVolumeFunc   func() ([]LogicalVolume, error)
	GetListLVMPhysicalVolumeFunc  func() ([]PhysicalVolume, error)
	GetgetSymLinksFunc            func(string) (string, error)

	fakeLogicalVolumeList = []LogicalVolume{
		{
			Name:     "pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
			FullName: "linuxlvmvg/pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
			UUID:     "AuUALf-Psow-usH5-IGDD-NeEC-TyAQ-ddBVeW",
			Size:     3221225472,
			Path:     "/dev/linuxlvmvg/pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
			DMPath:   "/dev/mapper/linuxlvmvg-pvc--3d584c82--a7ca--42bc--951d--29694fcc6637",
			VGName:   "linuxlvmvg",
			Active:   "active",
			Device:   "dm-1",
		},
		{
			Name:     "pvc-d264374b-2ab1-4ab7-b41a-48a2973dd0b8",
			FullName: "linuxlvmvg/pvc-d264374b-2ab1-4ab7-b41a-48a2973dd0b8",
			UUID:     "Ym6tU1-91Dc-3AE4-NpgQ-D35f-8rMr-B7F52s",
			Size:     10737418240,
			Path:     "/dev/linuxlvmvg/pvc-d264374b-2ab1-4ab7-b41a-48a2973dd0b8",
			DMPath:   "/dev/mapper/linuxlvmvg-pvc--d264374b--2ab1--4ab7--b41a--48a2973dd0b8",
			VGName:   "linuxlvmvg",
			Active:   "active",
			Device:   "dm-0",
		},
	}

	fakeLogicalVolume = LogicalVolume{

		Name:     "pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
		FullName: "linuxlvmvg/pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
		UUID:     "AuUALf-Psow-usH5-IGDD-NeEC-TyAQ-ddBVeW",
		Size:     3221225472,
		Path:     "/dev/linuxlvmvg/pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
		DMPath:   "/dev/mapper/linuxlvmvg-pvc--3d584c82--a7ca--42bc--951d--29694fcc6637",
		VGName:   "linuxlvmvg",
		Active:   "active",
	}

	fakeLVSOutput = "{\"report\":[{\"lv\":[{\"lv_uuid\":\"AuUALf-Psow-usH5-IGDD-NeEC-TyAQ-ddBVeW\",\"lv_name\":\"pvc-3d584c82-a7ca-42bc-951d-29694fcc6637\",\"lv_full_name\":\"linuxlvmvg/pvc-3d584c82-a7ca-42bc-951d-29694fcc6637\",\"lv_path\":\"/dev/linuxlvmvg/pvc-3d584c82-a7ca-42bc-951d-29694fcc6637\",\"lv_dm_path\":\"/dev/mapper/linuxlvmvg-pvc--3d584c82--a7ca--42bc--951d--29694fcc6637\",\"lv_parent\":\"\",\"lv_layout\":\"linear\",\"lv_role\":\"public\",\"lv_initial_image_sync\":\"\",\"lv_image_synced\":\"\",\"lv_merging\":\"\",\"lv_converting\":\"\",\"lv_allocation_policy\":\"inherit\",\"lv_allocation_locked\":\"\",\"lv_fixed_minor\":\"\",\"lv_skip_activation\":\"\",\"lv_when_full\":\"\",\"lv_active\":\"active\",\"lv_active_locally\":\"active locally\",\"lv_active_remotely\":\"\",\"lv_active_exclusively\":\"active exclusively\",\"lv_major\":\"-1\",\"lv_minor\":\"-1\",\"lv_read_ahead\":\"auto\",\"lv_size\":\"3221225472\",\"lv_metadata_size\":\"\",\"seg_count\":\"1\",\"origin\":\"\",\"origin_uuid\":\"\",\"origin_size\":\"\",\"lv_ancestors\":\"\",\"lv_full_ancestors\":\"\",\"lv_descendants\":\"\",\"lv_full_descendants\":\"\",\"raid_mismatch_count\":\"\",\"raid_sync_action\":\"\",\"raid_write_behind\":\"\",\"raid_min_recovery_rate\":\"\",\"raid_max_recovery_rate\":\"\",\"move_pv\":\"\",\"move_pv_uuid\":\"\",\"convert_lv\":\"\",\"convert_lv_uuid\":\"\",\"mirror_log\":\"\",\"mirror_log_uuid\":\"\",\"data_lv\":\"\",\"data_lv_uuid\":\"\",\"metadata_lv\":\"\",\"metadata_lv_uuid\":\"\",\"pool_lv\":\"\",\"pool_lv_uuid\":\"\",\"lv_tags\":\"\",\"lv_profile\":\"\",\"lv_lockargs\":\"\",\"lv_time\":\"2021-06-12 11:40:28 +0530\",\"lv_time_removed\":\"\",\"lv_host\":\"sumitworker1-virtual-machine\",\"lv_modules\":\"\",\"lv_historical\":\"\",\"lv_kernel_major\":\"253\",\"lv_kernel_minor\":\"0\",\"lv_kernel_read_ahead\":\"131072\",\"lv_permissions\":\"writeable\",\"lv_suspended\":\"\",\"lv_live_table\":\"live table present\",\"lv_inactive_table\":\"\",\"lv_device_open\":\"open\",\"data_percent\":\"\",\"snap_percent\":\"\",\"metadata_percent\":\"\",\"copy_percent\":\"\",\"sync_percent\":\"\",\"cache_total_blocks\":\"\",\"cache_used_blocks\":\"\",\"cache_dirty_blocks\":\"\",\"cache_read_hits\":\"\",\"cache_read_misses\":\"\",\"cache_write_hits\":\"\",\"cache_write_misses\":\"\",\"kernel_cache_settings\":\"\",\"kernel_cache_policy\":\"\",\"kernel_metadata_format\":\"\",\"lv_health_status\":\"\",\"kernel_discards\":\"\",\"lv_check_needed\":\"unknown\",\"lv_merge_failed\":\"unknown\",\"lv_snapshot_invalid\":\"unknown\",\"lv_attr\":\"-wi-ao----\",\"vg_name\":\"linuxlvmvg\"},{\"lv_uuid\":\"Ym6tU1-91Dc-3AE4-NpgQ-D35f-8rMr-B7F52s\",\"lv_name\":\"pvc-d264374b-2ab1-4ab7-b41a-48a2973dd0b8\",\"lv_full_name\":\"linuxlvmvg/pvc-d264374b-2ab1-4ab7-b41a-48a2973dd0b8\",\"lv_path\":\"/dev/linuxlvmvg/pvc-d264374b-2ab1-4ab7-b41a-48a2973dd0b8\",\"lv_dm_path\":\"/dev/mapper/linuxlvmvg-pvc--d264374b--2ab1--4ab7--b41a--48a2973dd0b8\",\"lv_parent\":\"\",\"lv_layout\":\"linear\",\"lv_role\":\"public\",\"lv_initial_image_sync\":\"\",\"lv_image_synced\":\"\",\"lv_merging\":\"\",\"lv_converting\":\"\",\"lv_allocation_policy\":\"inherit\",\"lv_allocation_locked\":\"\",\"lv_fixed_minor\":\"\",\"lv_skip_activation\":\"\",\"lv_when_full\":\"\",\"lv_active\":\"active\",\"lv_active_locally\":\"active locally\",\"lv_active_remotely\":\"\",\"lv_active_exclusively\":\"active exclusively\",\"lv_major\":\"-1\",\"lv_minor\":\"-1\",\"lv_read_ahead\":\"auto\",\"lv_size\":\"10737418240\",\"lv_metadata_size\":\"\",\"seg_count\":\"1\",\"origin\":\"\",\"origin_uuid\":\"\",\"origin_size\":\"\",\"lv_ancestors\":\"\",\"lv_full_ancestors\":\"\",\"lv_descendants\":\"\",\"lv_full_descendants\":\"\",\"raid_mismatch_count\":\"\",\"raid_sync_action\":\"\",\"raid_write_behind\":\"\",\"raid_min_recovery_rate\":\"\",\"raid_max_recovery_rate\":\"\",\"move_pv\":\"\",\"move_pv_uuid\":\"\",\"convert_lv\":\"\",\"convert_lv_uuid\":\"\",\"mirror_log\":\"\",\"mirror_log_uuid\":\"\",\"data_lv\":\"\",\"data_lv_uuid\":\"\",\"metadata_lv\":\"\",\"metadata_lv_uuid\":\"\",\"pool_lv\":\"\",\"pool_lv_uuid\":\"\",\"lv_tags\":\"\",\"lv_profile\":\"\",\"lv_lockargs\":\"\",\"lv_time\":\"2021-06-12 11:40:28 +0530\",\"lv_time_removed\":\"\",\"lv_host\":\"sumitworker1-virtual-machine\",\"lv_modules\":\"\",\"lv_historical\":\"\",\"lv_kernel_major\":\"253\",\"lv_kernel_minor\":\"1\",\"lv_kernel_read_ahead\":\"131072\",\"lv_permissions\":\"writeable\",\"lv_suspended\":\"\",\"lv_live_table\":\"live table present\",\"lv_inactive_table\":\"\",\"lv_device_open\":\"open\",\"data_percent\":\"\",\"snap_percent\":\"\",\"metadata_percent\":\"\",\"copy_percent\":\"\",\"sync_percent\":\"\",\"cache_total_blocks\":\"\",\"cache_used_blocks\":\"\",\"cache_dirty_blocks\":\"\",\"cache_read_hits\":\"\",\"cache_read_misses\":\"\",\"cache_write_hits\":\"\",\"cache_write_misses\":\"\",\"kernel_cache_settings\":\"\",\"kernel_cache_policy\":\"\",\"kernel_metadata_format\":\"\",\"lv_health_status\":\"\",\"kernel_discards\":\"\",\"lv_check_needed\":\"unknown\",\"lv_merge_failed\":\"unknown\",\"lv_snapshot_invalid\":\"unknown\",\"lv_attr\":\"-wi-ao----\",\"vg_name\":\"linuxlvmvg\"}]}]}"

	fakeVolumeGroupList = []VolumeGroup{
		{
			UUID:             "MBadZ4-QZYk-eGhp-kFd1-Dac1-PDMT-oG8TMz",
			Name:             "linuxlvmvg",
			Attr:             "wz--n-",
			Permission:       "writeable",
			AllocationPolicy: "normal",
			Format:           "lvm2",
			TotalSize:        26839351296,
			FreeSize:         12880707584,
			MaxLV:            0,
			MaxPV:            0,
			PVCount:          1,
			MissingPVCount:   0,
			LVCount:          2,
			SnapCount:        0,
		},
	}

	fakeVGSOutput = "{\"report\":[{\"vg\":[{\"vg_fmt\":\"lvm2\",\"vg_uuid\":\"MBadZ4-QZYk-eGhp-kFd1-Dac1-PDMT-oG8TMz\",\"vg_name\":\"linuxlvmvg\",\"vg_attr\":\"wz--n-\",\"vg_permissions\":\"writeable\",\"vg_extendable\":\"extendable\",\"vg_exported\":\"\",\"vg_partial\":\"\",\"vg_allocation_policy\":\"normal\",\"vg_clustered\":\"\",\"vg_size\":\"26839351296\",\"vg_free\":\"12880707584\",\"vg_sysid\":\"\",\"vg_systemid\":\"\",\"vg_lock_type\":\"\",\"vg_lock_args\":\"\",\"vg_extent_size\":\"4194304\",\"vg_extent_count\":\"6399\",\"vg_free_count\":\"3071\",\"max_lv\":\"0\",\"max_pv\":\"0\",\"pv_count\":\"1\",\"vg_missing_pv_count\":\"0\",\"lv_count\":\"2\",\"snap_count\":\"0\",\"vg_seqno\":\"21\",\"vg_tags\":\"\",\"vg_profile\":\"\",\"vg_mda_count\":\"1\",\"vg_mda_used_count\":\"1\",\"vg_mda_free\":\"0\",\"vg_mda_size\":\"1044480\",\"vg_mda_copies\":\"unmanaged\"}]}]}"

	fakePhysicalVolumeList = []PhysicalVolume{
		{

			UUID:      "TgirLZ-58x3-55Dn-A2wk-ZycJ-JUCh-aVeTf9",
			DevSize:   26843545600,
			Name:      "/dev/sdb1",
			TotalSize: 26839351296,
			FreeSize:  12880707584,
			UsedSize:  13958643712,
			Attr:      "a--",
			PVInUse:   "used",
		},
	}

	fakePVSOutput = "{\"report\":[{\"pv\":[{\"pv_fmt\":\"lvm2\",\"pv_uuid\":\"TgirLZ-58x3-55Dn-A2wk-ZycJ-JUCh-aVeTf9\",\"dev_size\":\"26843545600\",\"pv_name\":\"/dev/sdb1\",\"pv_major\":\"8\",\"pv_minor\":\"17\",\"pv_mda_free\":\"0\",\"pv_mda_size\":\"1044480\",\"pv_ext_vsn\":\"2\",\"pe_start\":\"1048576\",\"pv_size\":\"26839351296\",\"pv_free\":\"12880707584\",\"pv_used\":\"13958643712\",\"pv_attr\":\"a--\",\"pv_allocatable\":\"allocatable\",\"pv_exported\":\"\",\"pv_missing\":\"\",\"pv_pe_count\":\"6399\",\"pv_pe_alloc_count\":\"3328\",\"pv_tags\":\"\",\"pv_mda_count\":\"1\",\"pv_mda_used_count\":\"1\",\"pv_ba_start\":\"0\",\"pv_ba_size\":\"0\",\"pv_in_use\":\"used\",\"pv_duplicate\":\"\"}]}]}"
)

func (mock *MockLVMStruct) execLVMCommand(cmnd string, args ...string) ([]byte, error) {
	return GetexecLVMCommandFunc(cmnd, args...)
}
func (mock *MockLVMStruct) ReloadLVMMetadataCache() error {
	return GetReloadLVMMetadataCacheFunc()
}

func (moc *MockLVMStruct) getSymLinks(path string) (string, error) {
	return GetgetSymLinksFunc(path)
}

func (mock *MockLVMStruct) ListLVMVolumeGroup() ([]VolumeGroup, error) {
	return GetListLVMVolumeGroupFunc()
}

func (mock *MockLVMStruct) ListLVMLogicalVolume() ([]LogicalVolume, error) {
	return GetListLVMLogicalVolumeFunc()
}
func (mock *MockLVMStruct) ListLVMPhysicalVolume() ([]PhysicalVolume, error) {
	return GetListLVMPhysicalVolumeFunc()
}

func TestLVMstruct_ListLVMLogicalVolume(t *testing.T) {
	LVMClient = &MockLVMStruct{}
	tests := []struct {
		name    string
		want    []LogicalVolume
		wantErr bool
	}{
		{
			name:    "Test case to verify output of lvs command",
			want:    fakeLogicalVolumeList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvmC := &LVMstruct{}
			GetexecLVMCommandFunc = func(s string, s2 ...string) ([]byte, error) {
				return []byte(fakeLVSOutput), nil
			}
			GetReloadLVMMetadataCacheFunc = func() error {
				return nil
			}
			cnt := 0
			GetgetSymLinksFunc = func(s string) (string, error) {
				device := fakeLogicalVolumeList[cnt].Device
				cnt++
				return "/dev/" + device, nil
			}
			got, err := lvmC.ListLVMLogicalVolume()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListLVMLogicalVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListLVMLogicalVolume() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLVMstruct_ListLVMVolumeGroup(t *testing.T) {
	LVMClient = &MockLVMStruct{}
	tests := []struct {
		name    string
		want    []VolumeGroup
		wantErr bool
	}{
		{
			name:    "Test case to verify output of vgs command",
			want:    fakeVolumeGroupList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvmC := &LVMstruct{}
			GetexecLVMCommandFunc = func(s string, s2 ...string) ([]byte, error) {
				return []byte(fakeVGSOutput), nil
			}
			GetReloadLVMMetadataCacheFunc = func() error {
				return nil
			}
			got, err := lvmC.ListLVMVolumeGroup()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListLVMVolumeGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListLVMVolumeGroup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLVMstruct_ListLVMPhysicalVolume(t *testing.T) {
	LVMClient = &MockLVMStruct{}
	tests := []struct {
		name    string
		want    []PhysicalVolume
		wantErr bool
	}{
		{
			name:    "Test case to verify output of pvs command",
			want:    fakePhysicalVolumeList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lvmC := &LVMstruct{}
			GetexecLVMCommandFunc = func(s string, s2 ...string) ([]byte, error) {
				return []byte(fakePVSOutput), nil
			}
			GetReloadLVMMetadataCacheFunc = func() error {
				return nil
			}
			got, err := lvmC.ListLVMPhysicalVolume()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListLVMPhysicalVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListLVMPhysicalVolume() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLvDeviceName(t *testing.T) {
	type args struct {
		path string
	}
	LVMClient = &MockLVMStruct{}
	tests := []struct {
		name          string
		args          args
		want          string
		wantErr       bool
		symLinkOutput string
		symLinkError  error
	}{
		{
			name: "Test case with correct device name",
			args: args{
				"/dev/dm-0",
			},
			want:          "dm-0",
			wantErr:       false,
			symLinkOutput: "dm-0",
			symLinkError:  nil,
		},
		{
			name: "Test case with incorrect symlink",
			args: args{
				"fake-path",
			},
			want:          "",
			wantErr:       true,
			symLinkOutput: "",
			symLinkError:  errors.New("lvm: error in getting device name"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetgetSymLinksFunc = func(s string) (string, error) {
				return tt.symLinkOutput, tt.symLinkError
			}
			got, err := getLvDeviceName(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLvDeviceName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getLvDeviceName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLVMstruct_ReloadLVMMetadataCache(t *testing.T) {
	LVMClient = &MockLVMStruct{}
	tests := []struct {
		name              string
		want              error
		wantErr           bool
		execCommandOutput string
	}{
		{
			name:    "Test case with correct output",
			want:    nil,
			wantErr: false,
		},
		{
			name:    "Test case with incorrect output",
			want:    errors.New("Error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetexecLVMCommandFunc = func(s string, s2 ...string) ([]byte, error) {
				return []byte(""), tt.want
			}
			lvmC := &LVMstruct{}
			if err := lvmC.ReloadLVMMetadataCache(); (err != nil) != tt.wantErr {
				t.Errorf("ReloadLVMMetadataCache() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_parsePhysicalVolume(t *testing.T) {
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    PhysicalVolume
		wantErr bool
	}{
		{
			name: "Test case with successful parsing",
			args: args{
				map[string]string{
					"pv_fmt":    "lvm2",
					"pv_uuid":   "TgirLZ-58x3-55Dn-A2wk-ZycJ-JUCh-aVeTf9",
					"dev_size":  "26843545600",
					"pv_name":   "/dev/sdb1",
					"pe_start":  "1048576",
					"pv_size":   "26839351296",
					"pv_free":   "12880707584",
					"pv_used":   "13958643712",
					"pv_attr":   "a--",
					"pv_in_use": "used",
				},
			},
			want:    fakePhysicalVolumeList[0],
			wantErr: false,
		},
		{
			name: "Test case with unsuccessful parsing",
			args: args{
				map[string]string{
					"pv_fmt":    "lvm2",
					"pv_uuid":   "TgirLZ-58x3-55Dn-A2wk-ZycJ-JUCh-aVeTf9",
					"dev_size":  "invalid-format",
					"pv_name":   "/dev/sdb1",
					"pe_start":  "1048576",
					"pv_size":   "wrong-fmt",
					"pv_free":   "invalid-format",
					"pv_used":   "invalid-format",
					"pv_attr":   "a--",
					"pv_in_use": "used",
				},
			},
			want:    fakePhysicalVolumeList[0],
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parsePhysicalVolume(tt.args.m)
			if (err == nil) && tt.wantErr {
				t.Errorf("parsePhysicalVolume() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("parsePhysicalVolume() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseVolumeGroup(t *testing.T) {
	type args struct {
		m map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    VolumeGroup
		wantErr bool
	}{
		{
			name: "Test case with successful parsing",
			args: args{
				map[string]string{
					"vg_fmt":               "lvm2",
					"vg_uuid":              "MBadZ4-QZYk-eGhp-kFd1-Dac1-PDMT-oG8TMz",
					"vg_name":              "linuxlvmvg",
					"vg_attr":              "wz--n-",
					"vg_permissions":       "writeable",
					"vg_extendable":        "extendable",
					"vg_allocation_policy": "normal",
					"vg_size":              "26839351296",
					"vg_free":              "12880707584",
					"vg_extent_size":       "4194304",
					"vg_extent_count":      "6399",
					"vg_free_count":        "3071",
					"max_lv":               "0",
					"max_pv":               "0",
					"pv_count":             "1",
					"vg_missing_pv_count":  "0",
					"lv_count":             "2",
					"snap_count":           "0",
				},
			},
			want:    fakeVolumeGroupList[0],
			wantErr: false,
		},
		{
			name: "Test case with unsuccessful parsing",
			args: args{
				map[string]string{
					"vg_fmt":               "lvm2",
					"vg_uuid":              "MBadZ4-QZYk-eGhp-kFd1-Dac1-PDMT-oG8TMz",
					"vg_name":              "linuxlvmvg",
					"vg_attr":              "wz--n-",
					"vg_permissions":       "writeable",
					"vg_extendable":        "extendable",
					"vg_allocation_policy": "normal",
					"vg_clustered":         "",
					"vg_size":              "invalid-format",
					"vg_free":              "12880707584",
					"vg_extent_size":       "4194304",
					"vg_extent_count":      "6399",
					"vg_free_count":        "3071",
					"max_lv":               "0",
					"max_pv":               "0",
					"pv_count":             "1",
					"vg_missing_pv_count":  "0",
					"lv_count":             "2",
					"snap_count":           "0",
				},
			},
			want:    fakeVolumeGroupList[0],
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVolumeGroup(tt.args.m)
			if (err == nil) == tt.wantErr {
				t.Errorf("parseVolumeGroup() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("parseVolumeGroup() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseLogicalVolume(t *testing.T) {
	type args struct {
		m map[string]string
	}
	fakeLogicalVolumeList[0].Device = ""
	tests := []struct {
		name    string
		args    args
		want    LogicalVolume
		wantErr bool
	}{
		{
			name: "Test case with successful parsing",
			args: args{
				map[string]string{"lv_uuid": "AuUALf-Psow-usH5-IGDD-NeEC-TyAQ-ddBVeW",
					"lv_name":      "pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
					"lv_full_name": "linuxlvmvg/pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
					"lv_path":      "/dev/linuxlvmvg/pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
					"lv_dm_path":   "/dev/mapper/linuxlvmvg-pvc--3d584c82--a7ca--42bc--951d--29694fcc6637",
					"lv_size":      "3221225472", "vg_name": "linuxlvmvg", "lv_active": "active"},
			},
			want:    fakeLogicalVolume,
			wantErr: false,
		},
		{
			name: "Test case with unsuccessful parsing",
			args: args{
				map[string]string{"lv_uuid": "AuUALf-Psow-usH5-IGDD-NeEC-TyAQ-ddBVeW",
					"lv_name":      "pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
					"lv_full_name": "linuxlvmvg/pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
					"lv_path":      "/dev/linuxlvmvg/pvc-3d584c82-a7ca-42bc-951d-29694fcc6637",
					"lv_dm_path":   "/dev/mapper/linuxlvmvg-pvc--3d584c82--a7ca--42bc--951d--29694fcc6637",
					"lv_size":      "invalid-format", "vg_name": "linuxlvmvg", "lv_active": "active"},
			},
			want:    fakeLogicalVolume,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLogicalVolume(tt.args.m)
			if (err == nil) && tt.wantErr {
				t.Errorf("parseLogicalVolume() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("parseLogicalVolume() got = %v, want %v", got, tt.want)
			}
		})
	}
}
