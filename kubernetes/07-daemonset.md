[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----

## 目录

- [DaemonSet](#DaemonSet)
- [DaemonSet示例](#DaemonSet示例)


## DaemonSet

Linux有守护进程，k8s也有，由DaemonSet负责，DaemonSet负责在每个Node上部署一个相同功能的Pod，当然也可以不是全部的Node，通过selector部署到符合条件的Node上。

同Replicaset类似，DaemonSet也是Pod的Ownereference。

下图，给所有使用ssd的Node，部署上`ssd-monitor`。

![](http://img.lessisbetter.site/k8s-daemon.png)
*图片来自《Kubernetes in Action》*

DaemonSet可以：
- 在所有Node上部署会相同的Pod
- 监控集群节点，新加入的Node自动部署Pod
- 监控集群节点，确保被移除的Node删除Pod
- 监控每个Node上的Pod，退出会被拉起来

DaemonSet的Pod也是长时间运行的。

[↑top](#目录)


## DaemonSet示例

安装[官方文档](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/#required-fields)示例操作，本机Node是minikube，但却查不到daemonset，pod中也没有创建的迹象，从minikube的dashboard也没有发现daemonset。

是不是以为镜像下载不下来？

```sh
# 初次创建
➜  kubernetes git:(master) ✗ kubectl apply -f https://k8s.io/examples/controllers/daemonset.yaml
daemonset.apps/fluentd-elasticsearch created
# 
➜  kubernetes git:(master) ✗ kubectl get ds
No resources found in default namespace.
➜  kubernetes git:(master) ✗
➜  kubernetes git:(master) ✗ kubectl get node
NAME       STATUS   ROLES    AGE    VERSION
minikube   Ready    master   2d4h   v1.17.3
➜  kubernetes git:(master) ✗ kubectl get pod
No resources found in default namespace.

# 再次部署，但无创建
➜  kubernetes git:(master) ✗ kubectl apply -f daemonset.yaml
daemonset.apps/fluentd-elasticsearch unchanged
```

[↑top](#目录)
