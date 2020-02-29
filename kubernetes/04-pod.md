[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----

## 目录

- [Pod](#Pod)
- [Pod实现](#Pod实现)
- [容器设计模式Sidecar](#Sidecar)
    - [适配器容器](#适配器容器)
    - [代理容器](#代理容器)
    - [应用与日志收集](#应用与日志收集)

## Pod

Pod是k8s的资源和调度单位。

类比：
- K8s - 操作系统，比如Linux
- 容器 - 进程
- Pod - 进程组

具有超亲密关系的容器，应当放到一个Pod里面，不仅更容易实现上层应用的功能，更实现更高的效率。

![](http://img.lessisbetter.site/k8s-pod-concept.png)


[↑top](#目录)

### Pod实现

使用Linux namespace和cgroup实现容器间的隔离，一个Pod内的容器如何打破隔离，如何共享资源呢？

共享网络：
- 使用一个基础个容器打通Pod内容器间的网络通信，可以使用`localhost`访问其他容器
- 一个Pod只有1个IP，Pod内容器共享IP。
- Pod的声明周期与基础容器相同，要比Pod内容器长。

![](http://img.lessisbetter.site/k8s-pod-net.png)

共享存储：Node上目录挂在到Pod内每个容器上。

## SideCar

设计模式的本事是**重用与解耦**，编程语言的设计模式、架构的设计模式都是如此，容器设计模式也是如此。

容器的设计模式也使用SideCar，这让我想到了Service Mesh。

主容器负责核心功能，SideCar容器负责辅助工作，可以实现主容器功能的独立发布和能力重用。

[↑top](#目录)

### 适配器容器


适配器将主业务容器的接口转换为另外一个种格式。

![](http://img.lessisbetter.site/k8s-pod-adapter.png)


[↑top](#目录)

### 代理容器

屏蔽主容器也的业务集群，简化业务代码逻辑。

![](http://img.lessisbetter.site/k8s-pod-proxy.png)

[↑top](#目录)

### 应用与日志收集

主业务容器日志写磁盘，日志容器把日志从磁盘读出，然后上报进行收集。

![](http://img.lessisbetter.site/k8s-pod-log.png)

[↑top](#目录)