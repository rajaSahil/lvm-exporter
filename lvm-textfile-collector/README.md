## lvm textfile-collector script

lvm metric created by a Textfile Collector script. 

#### The following `metrics/labels` are collected:

```
lvm_lv_total_size{device="",lv_active="",lv_dm_path="",lv_full_name="",lv_name="",lv_path="",lv_uuid="",vg_name=""} TOTAL_SIZE
```

#### To publish container to Docker
* Run `make all` in the top directory. It will:
  * Build the docker image with the binary and will publish it in your docker repo.


#### To deploy:
- Run it as a sidecar container with node-exporter.
  - Please visit [examples](https://github.com/rajaSahil/lvm-exporter/tree/main/lvm-textfile-collector/examples) for `node-exporter` yaml and how to deploy it.