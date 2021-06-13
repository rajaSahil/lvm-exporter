package collector

import (
	"github.com/google/martian/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rajaSahil/lvm-exporter/pkg/lvm"
	"os/exec"
	"regexp"
	"strings"
)

type lvmCollector struct {
	vgFreeMetric *prometheus.Desc
	vgSizeMetric *prometheus.Desc

	lvSizeMetric *prometheus.Desc
}

// LVM Collector contains VG size and VG free in MB
func NewLvmCollector() *lvmCollector {
	return &lvmCollector{
		vgFreeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "free_size"),
			"Shows LVM VG free size",
			[]string{"vg_name"}, nil,
		),
		vgSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "vg", "total_size"),
			"Shows LVM VG total size",
			[]string{"vg_name"}, nil,
		),
		lvSizeMetric: prometheus.NewDesc(prometheus.BuildFQName("lvm", "lv", "total_size"),
			"Shows LVM LV total size",
			[]string{"lv_name", "lv_full_name", "lv_uuid", "lv_path", "lv_dm_path", "lv_active", "vg_name", "device"}, nil,
		),
	}
}

func (collector *lvmCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.vgFreeMetric
	ch <- collector.vgSizeMetric

	ch <- collector.lvSizeMetric
}

// LVM Collect, call OS command and set values
func (collector *lvmCollector) Collect(ch chan<- prometheus.Metric) {
	vgList, err := lvm.LVMClient.ListLVMVolumeGroup()
	if err != nil {
		log.Errorf("Error in getting the list of LVM volume groups:%v", err)
	}

	for _, vg := range vgList {
		ch <- prometheus.MustNewConstMetric(collector.vgFreeMetric, prometheus.GaugeValue, float64(vg.FreeSize), vg.Name)
		ch <- prometheus.MustNewConstMetric(collector.vgSizeMetric, prometheus.GaugeValue, float64(vg.TotalSize), vg.Name)
	}

	lvList, err := lvm.LVMClient.ListLVMLogicalVolume()
	if err != nil {
		log.Errorf("Error in getting the list of LVM logical volume:%v", err)
	}

	for _, lv := range lvList {
		ch <- prometheus.MustNewConstMetric(collector.lvSizeMetric, prometheus.GaugeValue, float64(lv.Size), lv.Name, lv.FullName, lv.UUID, lv.Path, lv.DMPath, lv.Active, lv.VGName, lv.Device)

	}
}

func getDeviceMap(vgMap *map[string]int64) map[string]string {
	deviceMap := make(map[string]string)
	for k, v := range *vgMap {
		if v == 0 {
			continue
		}
		out, err := exec.Command("/bin/ls", "-l", "/dev/"+k).Output()
		if err != nil {
			log.Errorf("Error: %v", err)
		}
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			space := regexp.MustCompile(`\s+`)
			s := space.ReplaceAllString(line, " ")
			values := strings.Split(s, " ")
			if len(values) == 11 {
				device := strings.Split(values[10], "/")
				deviceMap[values[8]] = device[1]
			}
		}

	}
	return deviceMap
}
