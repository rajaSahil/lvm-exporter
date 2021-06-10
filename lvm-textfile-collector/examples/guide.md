### To run lvm-exporter as a sidecar container with node-exporter:
```
Apply lvm-textfile-collector.yaml
kubectl apply -f lvm-textfile-collector.yaml
```

### To deploy it with `kube-prometheus-stack`

#### Steps:
1. Update `values.yaml`
    ```console
        ...
        prometheus-node-exporter:
          sidecars: 
            - image: sahil0071/lvm-textfile-collector
              name: monitor-lvm       
              securityContext:
                privileged: true
                runAsGroup: 0
                runAsNonRoot: false
                runAsUser: 0
              env:
                - name: TEXTFILE_PATH
                  value: /shared_vol
                - name: COLLECT_INTERVAL
                  value: "10"
              command:
                - /bin/bash
              args:
                - -c
                - ./textfile_collector -l
              volumeMounts:
                - mountPath: /shared_vol
                  name: textfile-collector
                - mountPath: /dev
                  name: dev-dir
          securityContext:
            fsGroup: 65534
            runAsGroup: 0
            runAsNonRoot: false
            runAsUser: 0
          extraHostVolumeMounts:
          - name: dev-dir
            readOnly: true
            mountPath: /dev
            hostPath: /dev
          - name: textfile-collector
            hostPath: /shared_vol
            mountPath: /shared_vol
            readOnly: true
            mountPropagation: None
          extraArgs:
            - --collector.textfile.directory=/shared_vol
            - --collector.filesystem.ignored-mount-points=^/(dev|proc|sys|var/lib/docker/.+)($|/)
            - --collector.filesystem.ignored-fs-types=^(tmpfs|autofs|binfmt_misc|cgroup|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|sysfs|tracefs)$
  
          ...
    ```
2. Run helm upgrade
    ```console
        helm upgrade -f values.yaml <RELEASE_NAME> <HELM_REPO_NAME>
    ```

### To deploy it with [openebs-monitoring](https://github.com/openebs/monitoring) to monitor lvm-localPV

#### Steps:
1. Update `values.yaml`
    ```console
      ...
      kube-prometheus-stack:
        prometheus-node-exporter:
          sidecars: 
            - image: sahil0071/lvm-textfile-collector
              name: monitor-lvm       
              securityContext:
                privileged: true
                runAsGroup: 0
                runAsNonRoot: false
                runAsUser: 0
              env:
                - name: TEXTFILE_PATH
                  value: /shared_vol
                - name: COLLECT_INTERVAL
                  value: "10"
              command:
                - /bin/bash
              args:
                - -c
                - ./textfile_collector -l
              volumeMounts:
                - mountPath: /shared_vol
                  name: textfile-collector
                - mountPath: /dev
                  name: dev-dir
          securityContext:
            fsGroup: 65534
            runAsGroup: 0
            runAsNonRoot: false
            runAsUser: 0
          extraHostVolumeMounts:
          - name: dev-dir
            readOnly: true
            mountPath: /dev
            hostPath: /dev
          - name: textfile-collector
            hostPath: /shared_vol
            mountPath: /shared_vol
            readOnly: true
            mountPropagation: None
          extraArgs:
            - --collector.textfile.directory=/shared_vol
            - --collector.filesystem.ignored-mount-points=^/(dev|proc|sys|var/lib/docker/.+)($|/)
            - --collector.filesystem.ignored-fs-types=^(tmpfs|autofs|binfmt_misc|cgroup|configfs|debugfs|devpts|devtmpfs|fusectl|hugetlbfs|mqueue|overlay|proc|procfs|pstore|rpc_pipefs|securityfs|sysfs|tracefs)$
  
      ...
    ```
2. Run helm upgrade
    ```console
        helm upgrade -f values.yaml <RELEASE_NAME> openebs-monitoring/openebs-monitoring
    ```