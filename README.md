# lvm-exporter

---

This repository provides code for a Prometheus metrics exporter
for LVM. This exporter can be run as a seperate daemonset or as a sidecar container.

This exporter makes use of
[prometheus-go-client](https://github.com/prometheus/client_golang), the official Go
client library for prometheus.

#### The following `metrics/labels` are being exported:

```
lvm_lv_total_size{device="",lv_active="",lv_dm_path="",lv_full_name="",lv_name="",lv_path="",lv_uuid="",vg_name=""} TOTAL_SIZE
lvm_vg_free_size{vg_name=""} FREE_SIZE
lvm_vg_total_size{vg_name=""} TOTAL_SIZE
```

#### To publish container to Docker
* Run `make all` in the top directory. It will:
  * Build the docker image with the binary and will publish it in your docker repo.


#### To deploy:
- Run it as a daemonset or as a sidecar container. 
  - Please visit [examples](https://github.com/rajaSahil/lvm-exporter/tree/main/examples) for daemonset yamls and how to deploy it.
  
## To run lvm-exporter as `textfile-collector` with node-exporter
- The textfile collector is similar to the Pushgateway, in that it allows exporting of statistics from batch jobs. 
It can also be used to export static metrics, such as what role a machine has. 
The Pushgateway should be used for service-level metrics. The textfile module is for metrics that are tied to a machine.

- To use it, set the **--collector.textfile.directory** flag on the node_exporter commandline. The collector will 
parse all files in that directory matching the glob ***.prom** using the text format. 

- To deploy lvm-textfile-collector with node-exporter:
  -  Please visit [lvm-textfile-collector](https://github.com/rajaSahil/lvm-exporter/tree/main/lvm-textfile-collector) to build image and use it with node-exporter
