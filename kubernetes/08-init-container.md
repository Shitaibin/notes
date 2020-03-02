[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----

### InitContainer

Initcontainer主要用来普通容器启动前的，初始化操作和前置条件检查，特点如下：

1. 先启动Initcontainer再启动普通容器
1. Initcontainer中可以定义多个容器，它们按定义的顺序依次启动（串行）
1. 普通容器在Initcontainer启动后，并发启动
1. Initcontainer执行后会退出，普通容器会持续运行或重启