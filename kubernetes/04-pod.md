[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----

## 目录

- [目录](#目录)
- [Pod](#pod)
  - [Pod与容器](#pod与容器)
  - [Pod实现](#pod实现)
  - [声明Pod](#声明pod)
- [SideCar](#sidecar)
  - [适配器容器](#适配器容器)
  - [代理容器](#代理容器)
  - [应用与日志收集](#应用与日志收集)
- [Pod与容器](#pod与容器-1)

## Pod

Pod是k8s的资源和调度单位，是K8S所推崇的**不可变基础设施**。

> 不可变基础设施则是另一种思路，部署完成以后，便成为一种只读状态，不可对其进行任何更改。如果需要更新或修改，就使用新的环境或服务器去替代旧的。不可变基础设施带来了更一致、更可靠、更可预测的设计理念，可以缓解或完全避免可变基础设施中遇到的各种常见问题。

![](https://github.com/kubernetes/website/blob/master/content/en/docs/tutorials/kubernetes-basics/public/images/module_03_pods.svg)

类比：
- K8s - 操作系统，比如Linux
- 容器 - 进程
- Pod - 进程组

具有超亲密关系的容器，应当放到一个Pod里面，不仅更容易实现上层应用的功能，更实现更高的效率，亲密关系有：
- 容器之间会发生文件交换
- 容器之间会存在本地通信，比如使用localhost
- 容器之间需要频繁的RPC
- 辅助容器：日志收集、监控数据采集、配置中心、路由及熔断

![](http://img.lessisbetter.site/k8s-pod-concept.png)
*图片来自阿里云原生公开课第4讲*



[↑top](#目录)

### Pod与容器

由于容器实际上是一个“单进程”的模型，这点非常重要。因为如果你在容器里启动多个进程，这将会带来很多麻烦。不仅它们的日志记录会混在一起，它们各自的生命周期也无法管理。毕竟只有一个进程的 PID 可以为 1，如果 PID 为 1 的进程这个时候挂了，或者说失败退出了，那么其他几个进程就会自然而然地成为“孤儿”，无法管理，也无法回收资源。

用一个 Pod 管理多个容器，既能够保持容器之间的隔离性，还能保证相关容器的环境一致性。使用粒度更小的容器，不仅可以使应用间的依赖解耦，还便于使用不同技术栈进行开发，同时还可以方便各个开发团队复用，减少重复造轮子。

[↑top](#目录)

### Pod实现

使用Linux namespace和cgroup实现容器间的隔离，一个Pod内的容器如何打破隔离，如何共享资源呢？

共享网络：
- 使用一个基础个容器打通Pod内容器间的网络通信，可以使用`localhost`访问其他容器
- 一个Pod只有1个IP，Pod内容器共享IP。
- Pod的声明周期与基础容器相同，要比Pod内容器长。

![](http://img.lessisbetter.site/k8s-pod-net.png)
*图片来自阿里云原生公开课第4讲*

共享存储：Node上目录挂在到Pod内每个容器上。

[↑top](#目录)

### 声明Pod

Pod包含3部分：
- metadata ：包含name、Namespace，声明当前Pod的名称和所在Namespace，还可以包含其他辅助内容
- spec ：期望Pod的“样子”，包含哪些容器，每个容器的期望是什么样子
- status ： Pod当前的状态

其中前2个可以通过yaml去声明。

## SideCar

设计模式的本事是**重用与解耦**，编程语言的设计模式、架构的设计模式都是如此，容器设计模式也是如此。

容器的设计模式也使用SideCar，这让我想到了Service Mesh。

主容器负责核心功能，SideCar容器负责辅助工作，可以实现主容器功能的独立发布和能力重用。

[↑top](#目录)

### 适配器容器


适配器将主业务容器的接口转换为另外一个种格式。

![](http://img.lessisbetter.site/k8s-pod-adapter.png)

> 本节内容图片来自论文 [Design patterns for container-based distributed systems](https://www.usenix.org/system/files/conference/hotcloud16/hotcloud16_burns.pdf) 。

[↑top](#目录)

### 代理容器

屏蔽主容器也的业务集群，简化业务代码逻辑。

![](http://img.lessisbetter.site/k8s-pod-proxy.png)

[↑top](#目录)

### 应用与日志收集

主业务容器日志写磁盘，日志容器把日志从磁盘读出，然后上报进行收集。

![](http://img.lessisbetter.site/k8s-pod-log.png)

[↑top](#目录)

## Pod与容器

连接到容器：

```
kubectl exec $pod_name -c $contianer_name -it /bin/bash
```

查看容器日志：

```
kubectl logs $pod_name -c $contianer_name
```