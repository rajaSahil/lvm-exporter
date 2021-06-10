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

#### To build the docker image:
```console
docker build .
docker tag <IMAGE_ID> REPO_NAME
docker push REPO_NAME
```

#### To deploy:
- Run it as a daemonset or as a sidecar container.