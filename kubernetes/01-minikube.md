玩转minikube
===============

minikube很好，但某些原因造成国内用起来比较慢，要各种挂代理、镜像加速。

## minikube原理

![](http://img.lessisbetter.site/2020-08-minikube.jpeg)

## 安装软件

1. 安装minikube，1分钟，如果提供的命令行下载不下来，就浏览器下载下来，放到增加可执行，然后放到bin目录即可：
https://yq.aliyun.com/articles/691500

1. centos安装virtualbox，2分钟安装完成:
https://wiki.centos.org/zh/HowTos/Virtualization/VirtualBox

3. 安装kubectl：
https://blog.csdn.net/yuanjunlai141/article/details/79469071


## 首次启动

启动命令
```
minikube start --image-mirror-country cn \
    --iso-url=https://kubernetes.oss-cn-hangzhou.aliyuncs.com/minikube/iso/minikube-v1.7.3.iso \
    --docker-env http_proxy=http://192.168.0.104:1087 \
    --docker-env https_proxy=http://192.168.0.104:1087 \
    --docker-env no_proxy=localhost,127.0.0.1,10.96.0.0/12,192.168.99.0/24,192.168.39.0/24 \
    --registry-mirror="https://a90tkz28.mirror.aliyuncs.com" \
    --image-repository="registry.cn-hangzhou.aliyuncs.com/google_containers" \
    --insecure-registry=192.168.9.8 \
    --kubernetes-version=v1.18.3
```

使用minikube可以查看帮助flag帮助信息：

- `--image-mirror-country`: 需要使用的镜像镜像的国家/地区代码。留空以使用全球代码。对于中国大陆用户，请将其设置为
cn
- `--docker-env`: 是通过环境变量向docker挂http代理，否则国内可能出现拉不到镜像的问题。挂代理还需要一个必要条件，在主机上使用SS开启代理。挂了代理可能也很难拉到，但不挂代理，几乎拉不下来镜像。
- `--registry-mirror`: 传递给 Docker 守护进程的注册表镜像。无效：--registry-mirror="https://a90tkz28.mirror.aliyuncs.com"
- `--image-repository` : 如果不能从gcr.io拉镜像，配置minikube中docker拉镜像的地方
- `--kubernetes-version`： 指定要部署的k8s版本

minikube内拉不到镜像的报错:

```
$ kubectl describe pod
  Type     Reason     Age                    From               Message
  ----     ------     ----                   ----               -------
  Warning  Failed     2m59s (x4 over 4m36s)  kubelet, minikube  Failed to pull image "kubeguide/redis-master": rpc error: code = Unknown desc = Error response from daemon: Get https://registry-1.docker.io/v2/: proxyconnect tcp: dial tcp 192.168.0.104:1087: connect: connection refused
```

启动日志：

```
$ minikube start --image-mirror-country cn \
    --iso-url=https://kubernetes.oss-cn-hangzhou.aliyuncs.com/minikube/iso/minikube-v1.7.3.iso \
    --docker-env http_proxy=http://192.168.0.104:1087 \
    --docker-env https_proxy=http://192.168.0.104:1087 \
    --docker-env no_proxy=localhost,127.0.0.1,10.96.0.0/12,192.168.99.0/24,192.168.39.0/24 \
    --registry-mirror="https://a90tkz28.mirror.aliyuncs.com" \
    --image-repository="registry.cn-hangzhou.aliyuncs.com/google_containers" \
    --insecure-registry=192.168.9.8
😄  Darwin 10.15.3 上的 minikube v1.12.3
✨  根据现有的配置文件使用 virtualbox 驱动程序
👍  Starting control plane node minikube in cluster minikube
🏃  Updating the running virtualbox "minikube" VM ...
🐳  正在 Docker 19.03.6 中准备 Kubernetes v1.18.3…
    ▪ env http_proxy=http://192.168.0.104:1087
    ▪ env https_proxy=http://192.168.0.104:1087
    ▪ env no_proxy=localhost,127.0.0.1,10.96.0.0/12,192.168.99.0/24,192.168.39.0/24
    > kubeadm.sha256: 65 B / 65 B [--------------------------] 100.00% ? p/s 0s
    > kubelet.sha256: 65 B / 65 B [--------------------------] 100.00% ? p/s 0s
    > kubeadm: 37.97 MiB / 37.97 MiB [--------------] 100.00% 320.45 MiB p/s 0s
    > kubelet: 108.04 MiB / 108.04 MiB [---------] 100.00% 514.43 KiB p/s 3m36s
🔎  Verifying Kubernetes components...
🌟  Enabled addons: default-storageclass, storage-provisioner
🏄  完成！kubectl 已经配置至 "minikube"
```

做哪些事？
1. 创建虚拟机"minikube"
2. 生成kubectl使用的配置文件，使用该配置连接集群：~/.kube/config
3. 在虚拟机里的容器上启动k8s

```
$ minikube ssh
                         _             _
            _         _ ( )           ( )
  ___ ___  (_)  ___  (_)| |/')  _   _ | |_      __
/' _ ` _ `\| |/' _ `\| || , <  ( ) ( )| '_`\  /'__`\
| ( ) ( ) || || ( ) || || |\`\ | (_) || |_) )(  ___/
(_) (_) (_)(_)(_) (_)(_)(_) (_)`\___/'(_,__/'`\____)

$
$ docker info
Client:
 Debug Mode: false

Server:
 Containers: 14
  Running: 14
  Paused: 0
  Stopped: 0
 Images: 10
 Server Version: 19.03.6
 Storage Driver: overlay2
  Backing Filesystem: extfs
  Supports d_type: true
  Native Overlay Diff: true
 Logging Driver: json-file
 Cgroup Driver: cgroupfs
 Plugins:
  Volume: local
  Network: bridge host ipvlan macvlan null overlay
  Log: awslogs fluentd gcplogs gelf journald json-file local logentries splunk syslog
 Swarm: inactive
 Runtimes: runc
 Default Runtime: runc
 Init Binary: docker-init
 containerd version: 35bd7a5f69c13e1563af8a93431411cd9ecf5021
 runc version: dc9208a3303feef5b3839f4323d9beb36df0a9dd
 init version: fec3683
 Security Options:
  seccomp
   Profile: default
 Kernel Version: 4.19.94
 Operating System: Buildroot 2019.02.9
 OSType: linux
 Architecture: x86_64
 CPUs: 2
 Total Memory: 3.754GiB
 Name: minikube
 ID: DSF4:HEQB:HTUU:OXRS:ZBWC:ESX4:WEST:UFDC:WAW5:5CDV:PITM:BEXZ
 Docker Root Dir: /var/lib/docker
 Debug Mode: false
 HTTP Proxy: http://192.168.0.104:1087
 HTTPS Proxy: http://192.168.0.104:1087
 No Proxy: localhost,127.0.0.1,10.96.0.0/12,192.168.99.0/24,192.168.39.0/24
 Registry: https://index.docker.io/v1/
 Labels:
  provider=virtualbox
 Experimental: false
 Insecure Registries:
  192.168.9.8
  10.96.0.0/12
  127.0.0.0/8
 Registry Mirrors:
  https://a90tkz28.mirror.aliyuncs.com/
 Live Restore Enabled: false
 Product License: Community Engine

$ exit
logout
```

Registry Mirrors对应的是阿里云镜像加速，HTTP proxy也配置上了，如果启动后，发现没有改变，需要删除过去创建的minikube，全部清理一遍。

## 常用minikube命令



- 集群状态： minikube status
- 暂停和恢复集群，不用的时候把它暂停掉，节约主机的CPU和内存： minikube pause， minikube unpause
- 停止集群： minikube stop
- 删除集群，遇到问题时，清理一波数据： minikube delete
- 查看集群IP，kubectl就是连这个IP： minikube ip
- 进入minikube虚拟机，整个k8s集群跑在这里面： minikube ssh

