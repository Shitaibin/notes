[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----


## 目录
- [目录](#目录)
- [Service简介](#service简介)
  - [ClusterIP：对集群内提供访问](#clusterip对集群内提供访问)
  - [NodePort：对集群外提供访问](#nodeport对集群外提供访问)
  - [LoadBalencer：对集群外提供访问](#loadbalencer对集群外提供访问)
  - [Ingress：对集群外提供访问](#ingress对集群外提供访问)
  - [Service架构](#service架构)
  - [Port & NodePort & TargetPort](#port--nodeport--targetport)
  - [Headless Service](#headless-service)


## Service简介

在[k8s介绍](./03-k8s.md#Service)中，介绍了k8s为什么需要Service，那么Service的原理是什么呢？怎么解决Pod与同Node内Pod、其他Node、集群外通信呢？

Service其实是一个代理，通过selector查找符合标签的Pod，Service是这些Pod的服务入口。

Service的yaml描述文件中指明了，service对外提供的端口，以及它要连接的Pod端口。

![](http://img.lessisbetter.site/k8s-service-real.png)

根据service是否被集群外访问分2类service，上图前端服务给集群外访问，后端服务给集群内访问。

集群内可以通过service的名字或者配置的IP访问service，service把请求转到Pod，如果有多个pod的情况，还会进行负载均衡，这个类型的service是ClusterIP。

集群外访问service，有3类方式：NodePort、LoadBalencer、InGress。

[↑top](#目录)



### ClusterIP：对集群内提供访问

这是Service的默认类型，可以为service指定ip，也可以不指定，当集群内访问`clusterip:port`时，service会把请求转给Pod。

Service所拥有的IP是虚拟IP，也称VIP。

![](http://img.lessisbetter.site/k8s-clusterip.png)

[↑top](#目录)

### NodePort：对集群外提供访问

把Node上的端口开放出来，让集群外可以访问，前提是节点的IP应当是公网可以访问的，NodePort上接到的请求都会转给Service，Service再转给Pod，由于Node和Pod直接多了一个Service，NodeA上接收的请求，有可能会被Service给NodeB上Pod去处理。

![](http://img.lessisbetter.site/k8s-nodeport.png)

[↑top](#目录)

### LoadBalencer：对集群外提供访问

云平台上通常有负载均衡服务，它相当于NodePort+负载均衡器。负载均衡器接收请求，转到隐式的NodePort，接下来流程如NodePort。

![](http://img.lessisbetter.site/k8s-loadbalencer.png)

[↑top](#目录)

### Ingress：对集群外提供访问

每个服务都需要自己的LoadBanlencer，当服务很多时，就需要创建很多LoadBalencer，有没有更好的方式解决这个问题？

Ingress可以解决这个问题，只需要一个公网ip就可以为很多服务提供外部接口，会根据主机名和路径会把请求转发到对应的service。

![](http://img.lessisbetter.site/k8s-ingress.png)

[↑top](#目录)


### Service架构

其架构，阿里云云原生公开课的这幅图更精确，相关组件的功能，也已经在图中标了出来。

![](http://img.lessisbetter.site/k8s-service-arch.png)

coredns负责service名字到service ip的转换。

如果service不设置ip，则通过service名字可以映射到pod的ip，直接访问pod。

[↑top](#目录)

### Port & NodePort & TargetPort

Port是service的端口。

NodePort是Node的端口。

TargetPort是Pod的端口。

[↑top](#目录)

### Headless Service

Service类型为ClusterIP，但无需分配IP地址给service。

headless service无法对请求路由，它只辅助提供生成从service名称到pod ip的映射，即域名记录。

部署[headless-svc.yaml](./examples/04-service/headless/headless-svr.yaml)后:

```
$ kubectl get svc,pod
NAME                             TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
service/headless-demo            ClusterIP      None             <none>        1234/TCP         7s

NAME                                   READY   STATUS        RESTARTS   AGE
pod/headless-demo-645dff5c6c-5qgsq     1/1     Running       0          7s
pod/headless-demo-645dff5c6c-clwj5     1/1     Running       0          7s
pod/headless-demo-645dff5c6c-w4628     1/1     Running       0          7s
```

可以看到service的CLUSTER-IP字段为None，service的描述信息可以看到已经关联的Endpoints，已经形成了service名字到Endpoints的映射。

```
$ kubectl describe svc headless-demo
Name:              headless-demo
Namespace:         default
Labels:            <none>
Annotations:       <none>
Selector:          name=headless-demo
Type:              ClusterIP
IP:                None
Port:              <unset>  1234/TCP
TargetPort:        1234/TCP
Endpoints:         10.32.0.14:1234,10.32.0.15:1234,10.32.0.16:1234
Session Affinity:  None
Events:            <none>
```

使用nslookup验证，确实能通过域名查询到3个对于的pod IP地址：

```
$ kubectl exec sleep -- nslookup headless-demo
Server:		10.96.0.10
Address:	10.96.0.10:53

Name:	headless-demo.default.svc.cluster.local
Address: 10.32.0.15
Name:	headless-demo.default.svc.cluster.local
Address: 10.32.0.16
Name:	headless-demo.default.svc.cluster.local
Address: 10.32.0.14
```


[↑top](#目录)