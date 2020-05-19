[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----


## 目录
- [目录](#%E7%9B%AE%E5%BD%95)
- [Service简介](#Service%E7%AE%80%E4%BB%8B)
  - [ClusterIP：对集群内提供访问](#ClusterIP%E5%AF%B9%E9%9B%86%E7%BE%A4%E5%86%85%E6%8F%90%E4%BE%9B%E8%AE%BF%E9%97%AE)
  - [NodePort：对集群外提供访问](#NodePort%E5%AF%B9%E9%9B%86%E7%BE%A4%E5%A4%96%E6%8F%90%E4%BE%9B%E8%AE%BF%E9%97%AE)
  - [LoadBalencer：对集群外提供访问](#LoadBalencer%E5%AF%B9%E9%9B%86%E7%BE%A4%E5%A4%96%E6%8F%90%E4%BE%9B%E8%AE%BF%E9%97%AE)
  - [InGress：对集群外提供访问](#InGress%E5%AF%B9%E9%9B%86%E7%BE%A4%E5%A4%96%E6%8F%90%E4%BE%9B%E8%AE%BF%E9%97%AE)
  - [Service架构](#Service%E6%9E%B6%E6%9E%84)
  - [Port & NodePort & TargetPort](#Port--NodePort--TargetPort)


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

### InGress：对集群外提供访问

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