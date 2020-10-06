[*Docker目录*](https://github.com/Shitaibin/notes/tree/master/docker#%E7%9B%AE%E5%BD%95)

-------------

## 目录

## AUFS

advanced multi-layered unification filesystem

特点：快速启动、高效利用存储和内存。

## Docker如何使用分层镜像

本机Docker信息，`Storage Driver: overlay2`可以看到使用的是overlay2作为存储驱动，而不是AUFS。

```
[/var/lib/docker]$ docker info
Client:
 Debug Mode: false

Server:
 Containers: 12
  Running: 6
  Paused: 0
  Stopped: 6
 Images: 23
 Server Version: 19.03.12
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
 containerd version: 7ad184331fa3e55e52b890ea95e65ba581ae3429
 runc version: dc9208a3303feef5b3839f4323d9beb36df0a9dd
 init version: fec3683
 Security Options:
  apparmor
  seccomp
   Profile: default
 Kernel Version: 4.15.0-112-generic
 Operating System: Ubuntu 18.04.3 LTS
 OSType: linux
 Architecture: x86_64
 CPUs: 4
 Total Memory: 7.789GiB
 Name: shitaibin-x
 ID: IX7L:SIYC:A6NP:2FHG:U6SD:5BK6:DRKT:O2SN:JWL2:EMKS:R3CQ:OFHE
 Docker Root Dir: /var/lib/docker
 Debug Mode: false
 Registry: https://index.docker.io/v1/
 Labels:
 Experimental: false
 Insecure Registries:
  127.0.0.0/8
 Registry Mirrors:
  https://a90tkz28.mirror.aliyuncs.com/
 Live Restore Enabled: false

WARNING: No swap limit support
```

```Dockerfile
[~]$ cat Dockerfile
From ubuntu:16.04
# Using Aliyun mirror
RUN echo "Hello wrold" > /tmp/newfile
```

Build得到镜像id `c0067aa3ef5d`：

```
[~]$ docker build -t changed-ubuntu .
Sending build context to Docker daemon  908.6MB
Step 1/2 : From ubuntu:16.04
 ---> 4b22027ede29
Step 2/2 : RUN echo "Hello wrold" > /tmp/newfile
 ---> Running in 8b401b1c0184
Removing intermediate container 8b401b1c0184
 ---> c0067aa3ef5d
Successfully built c0067aa3ef5d
Successfully tagged changed-ubuntu:latest
```

进入`/var/lib/docker`目录，搜索该镜像相关文件：

```
[/var/lib/docker]$ find . -name "c0067aa3ef5d*"
./image/overlay2/imagedb/metadata/sha256/c0067aa3ef5d766e1fe1dea666fa7fe1a52961b40ba1d291c7e9f6b9d6cc137b
./image/overlay2/imagedb/content/sha256/c0067aa3ef5d766e1fe1dea666fa7fe1a52961b40ba1d291c7e9f6b9d6cc137b
```

`./image/overlay2/imagedb/metadata/sha256/`目录下为镜像元数据：
- 镜像更新时间
- 该镜像的上一层镜像

```
[/var/lib/docker]$ cat  ./image/overlay2/imagedb/metadata/sha256/c0067aa3ef5d766e1fe1dea666fa7fe1a52961b40ba1d291c7e9f6b9d6cc137b/lastUpdated
2020-09-22T07:01:13.644666212Z#
[/var/lib/docker]$
[/var/lib/docker]$ cat  ./image/overlay2/imagedb/metadata/sha256/c0067aa3ef5d766e1fe1dea666fa7fe1a52961b40ba1d291c7e9f6b9d6cc137b/parent
sha256:4b22027ede299ea02d9d6236db8767e87b67392cf81535c18f7c202294a4a208#
```

`./image/overlay2/imagedb/content/sha256/`目录为镜像内容：

```
[/var/lib/docker]$ cat ./image/overlay2/imagedb/content/sha256/c0067aa3ef5d766e1fe1dea666fa7fe1a52961b40ba1d291c7e9f6b9d6cc137b | jq .
{
  "architecture": "amd64",
  "config": {
    "Hostname": "",
    "Domainname": "",
    "User": "",
    "AttachStdin": false,
    "AttachStdout": false,
    "AttachStderr": false,
    "Tty": false,
    "OpenStdin": false,
    "StdinOnce": false,
    "Env": [
      "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    ],
    "Cmd": [
      "/bin/bash"
    ],
    "ArgsEscaped": true,
    "Image": "sha256:4b22027ede299ea02d9d6236db8767e87b67392cf81535c18f7c202294a4a208",
    "Volumes": null,
    "WorkingDir": "",
    "Entrypoint": null,
    "OnBuild": null,
    "Labels": null
  },
  "container": "8b401b1c0184c4e4a194733db5840d95296b173ea99285ad0dbb0a100144560b",
  "container_config": {
    "Hostname": "",
    "Domainname": "",
    "User": "",
    "AttachStdin": false,
    "AttachStdout": false,
    "AttachStderr": false,
    "Tty": false,
    "OpenStdin": false,
    "StdinOnce": false,
    "Env": [
      "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
    ],
    "Cmd": [
      "|0",
      "/bin/sh",
      "-c",
      "echo \"Hello wrold\" > /tmp/newfile"
    ],
    "Image": "sha256:4b22027ede299ea02d9d6236db8767e87b67392cf81535c18f7c202294a4a208",
    "Volumes": null,
    "WorkingDir": "",
    "Entrypoint": null,
    "OnBuild": null,
    "Labels": null
  },
  "created": "2020-09-22T07:01:12.586900997Z",
  "docker_version": "19.03.12",
  "history": [
    {
      "created": "2020-08-19T21:16:18.060539271Z",
      "created_by": "/bin/sh -c #(nop) ADD file:144835a276ed2d8eaf6e893d5560444fe0d6a6f9b9bdadec1eb56e7bd9814427 in / "
    },
    {
      "created": "2020-08-19T21:16:19.582911032Z",
      "created_by": "/bin/sh -c rm -rf /var/lib/apt/lists/*"
    },
    {
      "created": "2020-08-19T21:16:21.047894836Z",
      "created_by": "/bin/sh -c set -xe \t\t&& echo '#!/bin/sh' > /usr/sbin/policy-rc.d \t&& echo 'exit 101' >> /usr/sbin/policy-rc.d \t&& chmod +x /usr/sbin/policy-rc.d \t\t&& dpkg-divert --local --rename --add /sbin/initctl \t&& cp -a /usr/sbin/policy-rc.d /sbin/initctl \t&& sed -i 's/^exit.*/exit 0/' /sbin/initctl \t\t&& echo 'force-unsafe-io' > /etc/dpkg/dpkg.cfg.d/docker-apt-speedup \t\t&& echo 'DPkg::Post-Invoke { \"rm -f /var/cache/apt/archives/*.deb /var/cache/apt/archives/partial/*.deb /var/cache/apt/*.bin || true\"; };' > /etc/apt/apt.conf.d/docker-clean \t&& echo 'APT::Update::Post-Invoke { \"rm -f /var/cache/apt/archives/*.deb /var/cache/apt/archives/partial/*.deb /var/cache/apt/*.bin || true\"; };' >> /etc/apt/apt.conf.d/docker-clean \t&& echo 'Dir::Cache::pkgcache \"\"; Dir::Cache::srcpkgcache \"\";' >> /etc/apt/apt.conf.d/docker-clean \t\t&& echo 'Acquire::Languages \"none\";' > /etc/apt/apt.conf.d/docker-no-languages \t\t&& echo 'Acquire::GzipIndexes \"true\"; Acquire::CompressionTypes::Order:: \"gz\";' > /etc/apt/apt.conf.d/docker-gzip-indexes \t\t&& echo 'Apt::AutoRemove::SuggestsImportant \"false\";' > /etc/apt/apt.conf.d/docker-autoremove-suggests"
    },
    {
      "created": "2020-08-19T21:16:22.502717186Z",
      "created_by": "/bin/sh -c mkdir -p /run/systemd && echo 'docker' > /run/systemd/container"
    },
    {
      "created": "2020-08-19T21:16:22.803326798Z",
      "created_by": "/bin/sh -c #(nop)  CMD [\"/bin/bash\"]",
      "empty_layer": true
    },
    {
      "created": "2020-09-22T07:01:12.586900997Z",
      "created_by": "|0 /bin/sh -c echo \"Hello wrold\" > /tmp/newfile"
    }
  ],
  "os": "linux",
  "rootfs": {
    "type": "layers",
    "diff_ids": [
      "sha256:e06660e80cf4fd425fd983f32d3682a9a6cd21728fb449d73e9272df6804bfcc",
      "sha256:41a253a417e6101be7476d5499cae21d90fde5caacb00e20efdcd17cc733962e",
      "sha256:87c1282613395a337e091a29f34bdd58f5f7200cc08ccf2613ecb30796a8348b",
      "sha256:dcc0cc99372e0cdd4d8f5e5308190caf01f03b981f34789bd825ad2faba41005",
      "sha256:1238119994ff24c29b4f934fbd21892443fc9c7d288d2db4d7e68f1d6fd62003"
    ]
  }
}
```