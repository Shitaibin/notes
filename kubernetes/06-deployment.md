[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----

## Deployment

Deloyment用来管理“部署”：创建、更新、回滚、删除部署。

- [Delpyment和Replicaset](#Delpyment和Replicaset)
- [Replicaset和Pod](#Replicaset和Pod)
- [Deloyment文件示例](#Deloyment文件示例)
- [更新Deloyment](#更新Deloyment)
- [回滚Deloyment](#回滚Deloyment)

### Delpyment和Replicaset

但Deloyment并不直接管理Pod，Deloyment管理的是Replicaset，Replicaset管理Pod，所以形成的OwnerReference如下：

![](http://img.lessisbetter.site/k8s-deloyment-replicaset.png)
*图片来自阿里云原生公开课第6讲*

所以Deloyment创建、更新、回滚的是Replicaset。

Deloyment和Replicaset是2个组件，Deloyment并不直接操控Replicaset，而是通过API Server。

![](http://img.lessisbetter.site/k8s-deloyment-arch.png)
*图片来自阿里云原生公开课第6讲*

![](http://img.lessisbetter.site/k8s-replicaset-arch.png)
*图片来自阿里云原生公开课第6讲*

```sh
# 创建了1个deployment
➜  kube kubectl get deployment
NAME               READY   UP-TO-DATE   AVAILABLE   AGE
nginx-deployment   5/5     5            5           16m
# deloyment有1个replicaset，id为nginx-deployment-54f57cf6bf
➜  kube kubectl get replicaset
NAME                          DESIRED   CURRENT   READY   AGE
nginx-deployment-54f57cf6bf   5         5         5       18m
```
[↑top](#Deployment)

### Replicaset和Pod

前面提到Pod是属于Replicaset的，而不是Deployment，解析一下Pod的名字:

- `nginx-deployment`：Deployment的名字。
- `54f57cf6bf`：Replicaset的ID。
- `44qwj`：Pod的ID，为随机数。

可以看到这3者的隶属关系。

```sh
➜  kube
# 当前replicaset有5个pod，pod名字的后5位为随机数
# pod是属于relicaset的，不是属于deployment
➜  kube kubectl get pod
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-54f57cf6bf-44qwj   1/1     Running   0          19m
nginx-deployment-54f57cf6bf-4q66c   1/1     Running   0          18m
nginx-deployment-54f57cf6bf-5rz8g   1/1     Running   0          19m
nginx-deployment-54f57cf6bf-nnptb   1/1     Running   0          18m
nginx-deployment-54f57cf6bf-wvfxl   1/1     Running   0          18m
```
[↑top](#Deployment)

### Deloyment文件示例

```yaml
apiVersion: apps/v1 # for versions before 1.9.0 use apps/v1beta2
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 5 # tells deployment to run 5 pods matching the template
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
```
文件来自 [k8s.io](https://k8s.io/examples/application/deployment.yaml) Deployment示例。

这里面包含：
- Deloyment信息
- Deployment的期望信息：
    - selector
    - replicas：pod的数量
    - template：pod的配置模板

[↑top](#Deployment)

### 更新Deloyment

Deloyment每次修改都会生成1个Replicaset，每个Replicaset都可以看做一个Deloyment部署实例版本，版本可以向后发展，也可以向前回滚。

![](http://img.lessisbetter.site/k8s-deployment-update.png)
*图片来自阿里云原生公开课第6讲*

```sh
# 修改deloyment配置
➜  kube kubectl set image deployment nginx-deployment nginx=nginx:1.9.1
deployment.apps/nginx-deployment image updated
# 重新配置pod的过程，产生了新的replicaset：nginx-deployment-56f8998dbc
➜  kube kubectl get pod
NAME                                READY   STATUS              RESTARTS   AGE
nginx-deployment-54f57cf6bf-44qwj   1/1     Running             0          21m
nginx-deployment-54f57cf6bf-4q66c   1/1     Running             0          20m
nginx-deployment-54f57cf6bf-5rz8g   1/1     Running             0          21m
nginx-deployment-54f57cf6bf-wvfxl   1/1     Running             0          20m
nginx-deployment-56f8998dbc-4cghd   0/1     ContainerCreating   0          63s
nginx-deployment-56f8998dbc-kklx5   0/1     ContainerCreating   0          63s
nginx-deployment-56f8998dbc-pr2wj   0/1     ContainerCreating   0          63s
# 修改deloyment后变成2个replicaset
# nginx-deployment-56f8998dbc有5个pod
# nginx-deployment-54f57cf6bf已经没有pod
➜  kube kubectl get replicaset
NAME                          DESIRED   CURRENT   READY   AGE
nginx-deployment-54f57cf6bf   0         0         0       43m
nginx-deployment-56f8998dbc   5         5         5       23m
➜  kube
➜  kube
# 从pod也可以看出，5个pod都属于56f8998dbc这个replicaset了
# 并且每个pod后面的随机数变了，说明pod是重新创建的
➜  kube kubectl get pod
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-56f8998dbc-4cghd   1/1     Running   0          21m
nginx-deployment-56f8998dbc-c8n6j   1/1     Running   0          16m
nginx-deployment-56f8998dbc-kklx5   1/1     Running   0          21m
nginx-deployment-56f8998dbc-ksgdm   1/1     Running   0          16m
nginx-deployment-56f8998dbc-pr2wj   1/1     Running   0          21m
```

[↑top](#Deployment)

### 回滚Deloyment

Deloyment下有多个Replicaset时，可以回滚到之前的某个Replicaset，回滚之后Replicaset会创建新的Pod，而不是使用之前该Replicaset拥有的Pod，因为那些Pod已经被删除。

```sh
# 回滚deloyment
➜  kube kubectl rollout undo deployment/nginx-deployment
deployment.apps/nginx-deployment rolled back
# 回滚过程中看到56f8998dbc还有1个pod，后面已被继续删除
# 54f57cf6bf已经有5个pod了
➜  kube kubectl get replicaset
NAME                          DESIRED   CURRENT   READY   AGE
nginx-deployment-54f57cf6bf   5         5         3       45m
nginx-deployment-56f8998dbc   1         1         1       25m
➜  kube
# 回滚完成，只有54f57cf6bf已经有5个pod
➜  kube kubectl get replicaset
NAME                          DESIRED   CURRENT   READY   AGE
nginx-deployment-54f57cf6bf   5         5         4       45m
nginx-deployment-56f8998dbc   0         0         0       25m
➜  kube
# pod也已经是使用nginx-deployment-54f57cf6bf
# pod后面的随机数与开始创建deloyment后初始的pod随机数并不相同，
# 说明这5个pod是重新创建的
➜  kube kubectl get pod
NAME                                READY   STATUS    RESTARTS   AGE
nginx-deployment-54f57cf6bf-7bxqv   1/1     Running   0          24s
nginx-deployment-54f57cf6bf-7qks8   1/1     Running   0          22s
nginx-deployment-54f57cf6bf-kkk8b   1/1     Running   0          23s
nginx-deployment-54f57cf6bf-l6r6j   1/1     Running   0          23s
nginx-deployment-54f57cf6bf-xvxpf   1/1     Running   0          22s
```

[↑top](#Deployment)