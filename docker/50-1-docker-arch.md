[*Docker目录*](https://github.com/Shitaibin/notes/tree/master/docker#%E7%9B%AE%E5%BD%95)

-------------

## 目录
- [目录](#目录)
- [Docker架构](#docker架构)
- [启动一个容器](#启动一个容器)
  - [参考资料](#参考资料)

## Docker架构

Docker三大件：
- 镜像：文件和文件夹的集合，封装了应用程序及其所有软件依赖的二进制数据
- 容器：容器运行着真正的应用进程。容器有初建、运行、停止、暂停和删除五种状态。
- 仓库：存储和分发 Docker 镜像

OCI为开放容器标准，所有的容器技术都实现了OCI，好方便容器编排软件对容器进行编排，分2部分：
- 运行时标准
- 镜像标准

Docker架构如下：

![](http://img.lessisbetter.site/2020-09-docker-arch.png)

Docker是CS架构。

Docker客户端：docker命令、docker sdk、调用RESTful API的浏览器等

Docker服务端：
- dockerd : 即Docker Daemon，负责接收和处理Docker客户端的请求
- containerd : 负责容器调度，通过containerd-shim管理runc。
- runc : 容器运行时，实现了OCI的运行时标准，它负责运行一个容器。

Docker包含的多个组件：

```
[~]$ ll /usr/bin/docker* /usr/bin/container* /usr/bin/runc /usr/bin/ctr
-rwxr-xr-x 1 root root  47M May  1 23:41 /usr/bin/containerd
-rwxr-xr-x 1 root root 5.9M May  1 23:41 /usr/bin/containerd-shim
-rwxr-xr-x 1 root root 7.4M May  1 23:41 /usr/bin/containerd-shim-runc-v1
-rwxr-xr-x 1 root root  25M May  1 23:41 /usr/bin/ctr
-rwxr-xr-x 1 root root  82M Jun 22 15:44 /usr/bin/docker
-rwxr-xr-x 1 root root 786K Jun 22 15:44 /usr/bin/docker-init
-rwxr-xr-x 1 root root 3.6M Jun 22 15:44 /usr/bin/docker-proxy
-rwxr-xr-x 1 root root  98M Jun 22 15:44 /usr/bin/dockerd
-rwxr-xr-x 1 root root 9.7M May  1 23:41 /usr/bin/runc
```

组件分3类：
1. docker
   - docker: docker客户端
   - dockerd: docker服务端，守护进程，即docker engine，dockerd是它可执行文件的缩写
   - docker-init: 容器内的进程为1号进程，当容器启动的进程无回收能力时，可以使用`--init`参数，让docker-init作为1号进程
   - docker-proxy: 负责容器的网络代理，当设置了端口映射`-p 8080:80`，主机80端口的请求，都会转发到容器的80端口，是通过iptabels实现的转发
2. [containerd](https://containerd.io/)，它是CNCF已毕业的项目，：
   - containerd: 负责容器管理，镜像拉取、存储和网络资源
   - containerd-shim: 利用shim实现containerd和真正容器进程的解耦。容器进程需要一个父进程来做诸如收集状态, 维持 stdin 等 fd 打开等工作. 而假如这个父进程就是 containerd, 那每次 containerd 挂掉或升级, 整个宿主机上所有的容器都得退出了。每个容器都有一个shim进程作为父进程，shim进程的声明周期和容器的声明周期相同。containerd-shim可以调用任何符合OCI的运行时，不仅仅是runc，还可以调用kata-runtime。通过这幅架构图，你更能理解containerd和shim的关系，可以把shim作为一个和具体容器运行时解耦的插件集合，它负责和各种具体的容器运行时交互。
   - ![containerd架构图](https://containerd.io/img/architecture.png)在Client那一层可以看到，kubelet通过CRI，dokcer、buildKit、ctr通过containerd client调用containerd。
   - ctr: 实际是containerd-ctr，是containerd的客户端，没有docker的时候，充当docker的部分功能

3. runc
   - runc: 是一个命令行工具，可以启动和运行1个容器，是 OCI (Open Container Initiative，开放容器标准)的标准实现。

> runc和containerd都是由最初的docker拆解出来的组件，使得开发者不需要通过docker engine也可以创建容器。runc是低层级的容器运行时，因为它只包含启动和运行容器，而containerd是更高层级的容器运行时，因为它还包含了镜像等其他功能。


## 启动一个容器

```
docker run -d busybox sleep 3600
```

通过进程树看 dockerd 、 containerd 、 runc 的关系：

```
[~]$ pstree -laA
systemd
  |-containerd
  |   |-containerd-shim -namespace moby -workdir /var/lib/containerd/io.containerd.runtime.v1.linux/moby/b78488a91722b3b40891944d74776f16e3a2ad70c11c47aaaf1443e6cd213c9c -address /run/containerd/containerd.sock -containerd-binary /usr/bin/containerd -runtime-root /var/run/docker/runtime-runc
  |   |   |-sleep 3600
  |   |   `-9*[{containerd-shim}]
  |   `-14*[{containerd}]
  |
  |-dockerd -H fd:// --containerd=/run/containerd/containerd.sock
  |   `-14*[{dockerd}]
```

dockerd启动时启动了containerd，containerd创建了containerd-shim，然后在其中创建了真正的进程`sleep 3600`。

### 参考资料

- [kubelet之cri演变史](https://zhuanlan.zhihu.com/p/87602649)，介绍了docker的演变，同时也就能清楚为何现在containerd和docker(engine)能平起平坐了，k8s这是要把docker(engine)直接踢出去了。