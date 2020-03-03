[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----

### Service

在[k8s介绍](./03-k8s.md#Service)中，介绍了k8s为什么需要Service，那么Service的原理是什么呢？怎么解决Pod与同Nod内Pod、其他Node、集群外通信呢？

Service其实是一个代理，通过selector找到要代理的Pod，而Service是这些Pod的服务入口。

Service的yaml描述文件中指明了，service对外提供的端口，以及它要连接的Pod端口。

同Node内其他Pod，可以通过service的名字或者IP访问service，service把请求转到Pod。

其他Node想访问service，service需要开放**NodePort**，形成Node端口到service端口的映射，由于Node的IP也是内部IP，集群外还无法访问。

设置**LoadBalancer**，提供公网IP和端口，形成公网IP和端口到Node端口的映射，公网可以通过LoadBalancer访问服务。

![](http://img.lessisbetter.site/k8s-service-nodeport.png)
*图片来自《Kubernetes in Action》*

其架构，阿里云云原生公开课的这幅图更精确，相关组件的功能，也已经在图中标了出来。

![](http://img.lessisbetter.site/k8s-service-arch.png)

coredns负责service名字到service ip的转换。

如果service不设置ip，则通过service名字可以映射到pod的ip，直接访问pod。