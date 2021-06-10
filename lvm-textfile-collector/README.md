## lvm textfile-collector script

lvm metric created by a Textfile Collector script. 

#### The following `metrics/labels` are collected:

```
lvm_lv_total_size{device="",lv_active="",lv_dm_path="",lv_full_name="",lv_name="",lv_path="",lv_uuid="",vg_name=""} TOTAL_SIZE
```

#### To build the docker image:
```console
docker build .
docker tag <IMAGE_ID> REPO_NAME
docker push REPO_NAME
```

#### To deploy:
- Run it as a sidecar container with node-exporter.