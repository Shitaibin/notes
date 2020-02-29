[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----

# 目录

- [lables](#lables)
- [anotations](#anotations)
- [ownereference](#ownereference)

### lables

KV类型，可以用来：
- 标识资源
- 筛选资源
- 组合资源
- 可使用selector查询，selector类型SQL

```sh
➜  kube kubectl get pods
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-54f57cf6bf-7bxqv   1/1     Running   0          117m
nginx-deployment-54f57cf6bf-7qks8   1/1     Running   0          117m
nginx-deployment-54f57cf6bf-kkk8b   1/1     Running   0          117m
nginx-deployment-54f57cf6bf-l6r6j   1/1     Running   0          117m
nginx-deployment-54f57cf6bf-xvxpf   1/1     Running   0          117m
➜  kube
# 展示pod标签
➜  kube kubectl get pods --show-labels
NAME                                READY   STATUS    RESTARTS   AGE    LABELS
nginx-deployment-54f57cf6bf-7bxqv   1/1     Running   0          117m   app=nginx,pod-template-hash=54f57cf6bf
nginx-deployment-54f57cf6bf-7qks8   1/1     Running   0          117m   app=nginx,pod-template-hash=54f57cf6bf
nginx-deployment-54f57cf6bf-kkk8b   1/1     Running   0          117m   app=nginx,pod-template-hash=54f57cf6bf
nginx-deployment-54f57cf6bf-l6r6j   1/1     Running   0          117m   app=nginx,pod-template-hash=54f57cf6bf
nginx-deployment-54f57cf6bf-xvxpf   1/1     Running   0          117m   app=nginx,pod-template-hash=54f57cf6bf
➜  kube
# 通过标签过滤
➜  kube kubectl get pods --show-labels -l app=nginx
NAME                                READY   STATUS    RESTARTS   AGE    LABELS
nginx-deployment-54f57cf6bf-7bxqv   1/1     Running   0          117m   app=nginx,pod-template-hash=54f57cf6bf
nginx-deployment-54f57cf6bf-7qks8   1/1     Running   0          117m   app=nginx,pod-template-hash=54f57cf6bf
nginx-deployment-54f57cf6bf-kkk8b   1/1     Running   0          117m   app=nginx,pod-template-hash=54f57cf6bf
nginx-deployment-54f57cf6bf-l6r6j   1/1     Running   0          117m   app=nginx,pod-template-hash=54f57cf6bf
nginx-deployment-54f57cf6bf-xvxpf   1/1     Running   0          117m   app=nginx,pod-template-hash=54f57cf6bf
```

[↑top](#目录)

### anotations

是资源注释，为KV类型，不用来标识资源，用来记录资源的非标识性资源，可以是结构化也可以是非结构化数据。

[↑top](#目录)

### ownereference

用来表示层级的资源关系，Replicaset是Pod是ownereference，Deployment是Replicaset的ownereference。

![](http://img.lessisbetter.site/k8s-deloyment-replicaset.png)
*图片来自阿里云原生公开课第6讲*

[↑top](#目录)