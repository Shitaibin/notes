


本节展示了本地存储的PVC、PV、SC用法。

验证1：sc的绑定模式设置了`WaitForFirstConsumer`，只有当pvc被Pod使用时，pvc和pv才会绑定。

以下记录为验证：先创建pv、pvc、sc，而不创建deployment，pvc的描述信息中显示了sc的名字，但未显示pv，说明pvc没有绑定pv。创建deployment后，查看pv、pvc状态都为已绑定。

```
$ kubectl apply -f pvc.yaml -f storage-class.yaml
persistentvolumeclaim/local-pvc created
storageclass.storage.k8s.io/local-sc created
$
$ kubectl get pvc
NAME            STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS   AGE
local-pvc       Pending                                      local-sc       5s
pd-basic-pd-0   Pending                                                     131m
$ kubectl get sc
NAME       PROVISIONER                    RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
local-sc   kubernetes.io/no-provisioner   Delete          WaitForFirstConsumer   false                  9s
$ kubectl describe pvc local-pv
Name:          local-pvc
Namespace:     default
StorageClass:  local-sc
Status:        Pending
Volume:
Labels:        <none>
Annotations:   <none>
Finalizers:    [kubernetes.io/pvc-protection]
Capacity:
Access Modes:
VolumeMode:    Filesystem
Mounted By:    <none>
Events:
  Type     Reason                Age   From                         Message
  ----     ------                ----  ----                         -------
  Warning  ProvisioningFailed    18s   persistentvolume-controller  storageclass.storage.k8s.io "local-sc" not found
  Normal   WaitForFirstConsumer  11s   persistentvolume-controller  waiting for first consumer to be created before binding

$
$ kubectl apply -f pv.yaml
persistentvolume/local-pv created
$
$ kubectl get pvc
NAME            STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS   AGE
local-pvc       Pending                                      local-sc       86s
pd-basic-pd-0   Pending                                                     132m
$ kubectl get pv
NAME       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM   STORAGECLASS   REASON   AGE
local-pv   5Gi        RWO            Recycle          Available           local-sc                8s
$ kubectl describe pvc local-pv
Name:          local-pvc
Namespace:     default
StorageClass:  local-sc
Status:        Pending
Volume:
Labels:        <none>
Annotations:   <none>
Finalizers:    [kubernetes.io/pvc-protection]
Capacity:
Access Modes:
VolumeMode:    Filesystem
Mounted By:    <none>
Events:
  Type     Reason                Age               From                         Message
  ----     ------                ----              ----                         -------
  Warning  ProvisioningFailed    103s              persistentvolume-controller  storageclass.storage.k8s.io "local-sc" not found
  Normal   WaitForFirstConsumer  6s (x7 over 96s)  persistentvolume-controller  waiting for first consumer to be created before binding

$
$ kubectl apply -f deployment.yaml
deployment.apps/busybox created
$ kubectl get pvc
NAME            STATUS    VOLUME     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
local-pvc       Bound     local-pv   5Gi        RWO            local-sc       2m16s
pd-basic-pd-0   Pending                                                       133m
$ kubectl describe pvc local-pv
Name:          local-pvc
Namespace:     default
StorageClass:  local-sc
Status:        Bound
Volume:        local-pv
Labels:        <none>
Annotations:   pv.kubernetes.io/bind-completed: yes
               pv.kubernetes.io/bound-by-controller: yes
Finalizers:    [kubernetes.io/pvc-protection]
Capacity:      5Gi
Access Modes:  RWO
VolumeMode:    Filesystem
Mounted By:    busybox-6f98c8d7fb-gfvp4
Events:
  Type     Reason                Age                  From                         Message
  ----     ------                ----                 ----                         -------
  Warning  ProvisioningFailed    2m20s                persistentvolume-controller  storageclass.storage.k8s.io "local-sc" not found
  Normal   WaitForFirstConsumer  13s (x9 over 2m13s)  persistentvolume-controller  waiting for first consumer to be created before binding

$ kubectl get pv
NAME       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM               STORAGECLASS   REASON   AGE
local-pv   5Gi        RWO            Recycle          Bound    default/local-pvc   local-sc                70s
```

验证2：StorageClass为本地存储（`provisioner: kubernetes.io/no-provisioner`）的情况下，pvc找到sc后，sc不会自动创建pv，需要手动创建pv。验证如下：创建pvc、sc、deployment，而不创建pv，发现pod、pvc都报没有可用pv的提示，并且pvc的描述信息显示sc的名字，而没有pv的名字，创建pv后，pv、pvc的状态都为已绑定，pod可以被调度然后运行。

```
$ kubectl apply -f pvc.yaml -f storage-class.yaml
persistentvolumeclaim/local-pvc created
storageclass.storage.k8s.io/local-sc created
$ kubectl apply -f deployment.yaml
deployment.apps/busybox created
$ kubectl get pvc
NAME            STATUS    VOLUME   CAPACITY   ACCESS MODES   STORAGECLASS   AGE
local-pvc       Pending                                      local-sc       10s
pd-basic-pd-0   Pending                                                     138m
$ kubectl get pod
NAME                               READY   STATUS    RESTARTS   AGE
basic-discovery-5bcf68669b-p72g9   1/1     Running   0          132m
basic-pd-0                         0/1     Pending   0          132m
busybox-6f98c8d7fb-f2fv5           0/1     Pending   0          13s
mysql-0                            1/1     Running   0          6d7h
$ kubectl describe pod busybox-6f98c8d7fb-f2fv5|tail
    SecretName:  default-token-w22r7
    Optional:    false
QoS Class:       BestEffort
Node-Selectors:  <none>
Tolerations:     node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                 node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type     Reason            Age                From               Message
  ----     ------            ----               ----               -------
  Warning  FailedScheduling  23s (x2 over 23s)  default-scheduler  0/1 nodes are available: 1 node(s) didn't find available persistent volumes to bind.
$
$ kubectl describe pvc local-pv
Name:          local-pvc
Namespace:     default
StorageClass:  local-sc
Status:        Pending
Volume:
Labels:        <none>
Annotations:   <none>
Finalizers:    [kubernetes.io/pvc-protection]
Capacity:
Access Modes:
VolumeMode:    Filesystem
Mounted By:    busybox-6f98c8d7fb-f2fv5
Events:
  Type     Reason               Age               From                         Message
  ----     ------               ----              ----                         -------
  Warning  ProvisioningFailed   38s               persistentvolume-controller  storageclass.storage.k8s.io "local-sc" not found
  Normal   WaitForPodScheduled  2s (x3 over 32s)  persistentvolume-controller  waiting for pod busybox-6f98c8d7fb-f2fv5 to be scheduled

$ kubectl apply -f pv.yaml
persistentvolume/local-pv created
$
$ kubectl get pvc,pv
NAME                                  STATUS    VOLUME     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
persistentvolumeclaim/local-pvc       Bound     local-pv   5Gi        RWO            local-sc       64s
persistentvolumeclaim/pd-basic-pd-0   Pending                                                       139m

NAME                        CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM               STORAGECLASS   REASON   AGE
persistentvolume/local-pv   5Gi        RWO            Recycle          Bound    default/local-pvc   local-sc                16s
```

验证3：如果pvc被pod所使用，删除pvc的操作会被阻塞，pvc的状态始终为`Terminating`，直到pod被删除后，pvc的删除操作才完成。

验证4：pvc依赖sc时，在没有创建sc的情况下，即使手动创建了满足pvc条件的pv，pvc和pv也无法绑定。验证如下：创建pv、pvc后，未显示bound，因为pvc使用的sc。创建sc后pv、pvc为bound。

```
$ kubectl apply -f pv.yaml -f pvc.yaml
persistentvolume/local-pv created
persistentvolumeclaim/local-pvc created
$
$ kubectl get pvc
NAME            STATUS    VOLUME     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
local-pvc       Bound     local-pv   5Gi        RWO            local-sc       6s
pd-basic-pd-0   Pending                                                       124m
$ kubectl get pvc
NAME            STATUS    VOLUME     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
local-pvc       Bound     local-pv   5Gi        RWO            local-sc       7s
pd-basic-pd-0   Pending                                                       124m

$ kubectl apply -f storage-class.yaml
storageclass.storage.k8s.io/local-sc created
$ kubectl get pvc
NAME            STATUS    VOLUME     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
local-pvc       Bound     local-pv   5Gi        RWO            local-sc       26s
pd-basic-pd-0   Pending                                                       125m
$ kubectl get pvc
NAME            STATUS    VOLUME     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
local-pvc       Bound     local-pv   5Gi        RWO            local-sc       31s
pd-basic-pd-0   Pending                                                       125m
$ kubectl get pv
NAME       CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS   CLAIM               STORAGECLASS   REASON   AGE
local-pv   5Gi        RWO            Recycle          Bound    default/local-pvc   local-sc                33s
$ kubectl get sc
NAME       PROVISIONER                    RECLAIMPOLICY   VOLUMEBINDINGMODE      ALLOWVOLUMEEXPANSION   AGE
local-sc   kubernetes.io/no-provisioner   Delete          WaitForFirstConsumer   false                  18s
$ kubectl describe sc local-sc
Name:            local-sc
IsDefaultClass:  No
Annotations:     kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"storage.k8s.io/v1","kind":"StorageClass","metadata":{"annotations":{},"name":"local-sc"},"provisioner":"kubernetes.io/no-provisioner","volumeBindingMode":"WaitForFirstConsumer"}

Provisioner:           kubernetes.io/no-provisioner
Parameters:            <none>
AllowVolumeExpansion:  <unset>
MountOptions:          <none>
ReclaimPolicy:         Delete
VolumeBindingMode:     WaitForFirstConsumer
Events:                <none>
$
$
$ kubectl describe pvc local-pvc
Name:          local-pvc
Namespace:     default
StorageClass:  local-sc  # 所使用的sc，如果为空，说明未指定sc，会使用默认sc
Status:        Bound
Volume:        local-pv  # 所绑定的pv
Labels:        <none>
Annotations:   pv.kubernetes.io/bind-completed: yes
               pv.kubernetes.io/bound-by-controller: yes
Finalizers:    [kubernetes.io/pvc-protection]
Capacity:      5Gi
Access Modes:  RWO
VolumeMode:    Filesystem
Mounted By:    <none>
Events:        <none>
```


