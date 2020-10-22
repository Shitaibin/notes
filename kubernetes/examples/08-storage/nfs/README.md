
1. 以Service提供NFS服务。 
busybox-nfs -> PVC -> StorageClass(自动创建PV) -> NFS service

NFS servcie使用Deployment，Deployment可以再使用volume，或依赖PVC，依赖公有云的存储。
NFS service -> NFS deployment -> volume
                              -> PVC -> 公有云

操作步骤：

1. 创建ServiceAccount给NFS Pod、Busybox Pod使用: `kubectl create serviceaccount nfs-provisioner`。
2. 创建NFS server的deployment、service。
3. 创建busybox使用的storageclass、pvc、deployment。