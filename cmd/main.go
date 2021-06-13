package main

import (
	"fmt"
	"github.com/go-kit/kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/rajaSahil/lvm-exporter/pkg/collector"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"os"
)

var (
	enableLvmExporter      = kingpin.Flag("web.enable-lvm-exporter", "Enable lvm exporter").Default("true").Bool()
	listenAddress          = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface.").Default(":9880").String()
	metricsPath            = kingpin.Flag("web.metrics-path", "Path under which to expose metrics.").Default("/metrics").String()
	disableExporterMetrics = kingpin.Flag("web.disable-exporter-metrics", "Exclude metrics about the exporter itself (promhttp_*, process_*, go_*).").Default("true").Bool()
)

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("node_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)
	fmt.Println(level.Info(logger).Log("msg", "Starting lvm-exporter", "version", version.Info()))
	fmt.Println(level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext()))

	registry := prometheus.NewRegistry()

	if !*disableExporterMetrics {
		registry.MustRegister(
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			collectors.NewGoCollector(),
		)
	}
	registry.MustRegister(version.NewCollector("lvm_exporter"))

	lvmExporter := collector.NewLvmCollector()

	registry.MustRegister(lvmExporter)

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`<html>
			<head><title>LVM Exporter</title></head>
			<body>
			<h1>LVM Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	_ = level.Info(logger).Log("msg", "Listening on", "address", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		_ = level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
