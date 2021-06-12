package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os/exec"
	"regexp"
	"strconv"
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
	out, err := exec.Command("/sbin/vgs", "--units", "M", "--separator", ",", "-o", "vg_name,vg_free,vg_size,lv_count", "--noheadings").Output()
	vgMap := make(map[string]int64)

	if err != nil {
		log.Print(err)
	}
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		values := strings.Split(line, ",")
		if len(values) == 4 {
			free_size, err := strconv.ParseFloat(strings.Trim(values[1], "M"), 64)
			if err != nil {
				log.Print(err)
			} else {
				total_size, err := strconv.ParseFloat(strings.Trim(values[2], "M"), 64)
				if err != nil {
					log.Print(err)
				} else {
					vg_name := strings.Trim(values[0], " ")
					lv_count := strings.Trim(values[3], " ")
					vgMap[vg_name], _ = strconv.ParseInt(lv_count, 0, 64)
					ch <- prometheus.MustNewConstMetric(collector.vgFreeMetric, prometheus.GaugeValue, free_size, vg_name)
					ch <- prometheus.MustNewConstMetric(collector.vgSizeMetric, prometheus.GaugeValue, total_size, vg_name)
				}
			}
		}
	}
	deviceMap := getDeviceMap(&vgMap)

	out, err = exec.Command("/sbin/lvs", "--units", "M", "--separator", ",", "-o", "lv_name,lv_full_name,lv_uuid,lv_path,lv_dm_path,lv_size,lv_active,vg_name", "--noheadings").Output()

	if err != nil {
		log.Println(err)
	}

	lines = strings.Split(string(out), "\n")

	for _, line := range lines {
		values := strings.Split(line, ",")
		if len(values) == 8 {
			total_size, err := strconv.ParseFloat(strings.Trim(values[5], "M"), 64)
			if err != nil {
				log.Print(err)
			} else {
				lv_name := strings.Trim(values[0], " ")
				lv_full_name := strings.Trim(values[1], " ")
				lv_uuid := strings.Trim(values[2], " ")
				lv_path := strings.Trim(values[3], " ")
				lv_dm_path := strings.Trim(values[4], " ")
				lv_active := strings.Trim(values[6], " ")
				vg_name := strings.Trim(values[7], " ")
				ch <- prometheus.MustNewConstMetric(collector.lvSizeMetric, prometheus.GaugeValue, total_size, lv_name, lv_full_name, lv_uuid, lv_path, lv_dm_path, lv_active, vg_name, deviceMap[lv_name])
			}

		}
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
			log.Println(err)
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
