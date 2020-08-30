


## cpu 

### 默认的情况

stress为基于ubuntu:16.04安装stress做出来的镜像，启动容器不做任何cpu限制。

```
[/sys/fs/cgroup/cpu]$ docker run --rm -it  -m 128m stress:16.04
root@5fad38726740:/# stress --vm-bytes 100m --vm-keep -m 1 -c 2
stress: info: [12] dispatching hogs: 2 cpu, 0 io, 1 vm, 0 hdd
```

查看 `cfs_period_us` 和 `cfs_quota_us` 的默认值：

```
[/sys/fs/cgroup/cpu,cpuacct/system.slice/docker-5fad38726740b90b93c06972fe4a9f11391a38aaeb3e922f10c3269fa32e1873.scope]$ cat cpu.cfs_period_us
100000
[/sys/fs/cgroup/cpu,cpuacct/system.slice/docker-5fad38726740b90b93c06972fe4a9f11391a38aaeb3e922f10c3269fa32e1873.scope]$ cat cpu.cfs_quota_us
-1
```

查看主机CPU利用率，为3个stress进程，每1个都100%，它们属于同一个cgroup：

```
  PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
 5616 root      20   0  109872 102336     36 R 100.0  1.3   0:07.46 stress
 5617 root      20   0    7468     88      0 R 100.0  0.0   0:07.45 stress
 5615 root      20   0    7468     88      0 R 100.0  0.0   0:07.45 stress
```


### 设置的情况

设置5000，在主机上只能占5%的cpu，开stress，分配到另外2个核。

[/sys/fs/cgroup/cpu]$ docker run --rm -it  -m 128m --cpu-quota=5000 stress:16.04
root@7e79005d7ca1:/#
root@7e79005d7ca1:/#
root@7e79005d7ca1:/# stress --vm-bytes 100m --vm-keep -m 1 -c 2
stress: info: [13] dispatching hogs: 2 cpu, 0 io, 1 vm, 0 hdd

查看 `cfs_period_us` 和 `cfs_quota_us` 的设置，5000/100000 = 5% ， 即限制该容器的CPU使用率不得超过5%。

[/sys/fs/cgroup/cpu,cpuacct/system.slice/docker-7e79005d7ca1b338d870d3dc79af3f1cd38ace195ebd685a09575f6acee36a07.scope]$ cat cpu.cfs_quota_us
5000
[/sys/fs/cgroup/cpu,cpuacct/system.slice/docker-7e79005d7ca1b338d870d3dc79af3f1cd38ace195ebd685a09575f6acee36a07.scope]$ cat cpu.cfs_period_us
100000

可以看到3个进程总cpu使用率5.1%。

```
  PID USER      PR  NI    VIRT    RES    SHR S  %CPU %MEM     TIME+ COMMAND
 5411 root      20   0    7468     92      0 R   1.7  0.0   0:00.53 stress
 5412 root      20   0  109872 102500     36 R   1.7  1.3   0:00.30 stress
 5413 root      20   0    7468     92      0 R   1.7  0.0   0:00.35 stress
```

## cpuacct



要同时输出，不然可能有时间差。

```
[/sys/fs/cgroup/cpu,cpuacct]$ cat cpuacct.*
user 20244450 // cpuacct.stat
system 52361  // cpuacct.stat
204310768947624  // cpuacct.usage
61143521333219 32616883199042 73804985004267 36745379411096 // cpuacct.usage_percpu
[/sys/fs/cgroup/cpu,cpuacct]$ python
Python 2.7.5 (default, Apr 11 2018, 07:36:10)
[GCC 4.8.5 20150623 (Red Hat 4.8.5-28)] on linux2
Type "help", "copyright", "credits" or "license" for more information.
>>> sum = 61143521333219+32616883199042+73804985004267+36745379411096
>>> dist = sum - 204310768947624
>>> dist
0
>>> sum
204310768947624
>>> sum2 = 20244450+52361
>>> sum2
20296811
```

## cpuset


## docker和cgroup查看cpu


6. 查看cgroup中的cpu信息

之前的容器样例限制了内存，没有限制cpu，我们启动一个限制cpu的容器。

```
docker run --rm -itd -c 2 -m 128m ubuntu:16.04
30a1da397ec45034cca1582c6a2be9095693a02196a5e40a531ffd4abda1f33d
```

然后另外开1个进程，登录该容器，查看容器中的进程：

```
docker exec -it 30a1da397ec4 /bin/bash
root@30a1da397ec4:/# ps -A
  PID TTY          TIME CMD
    1 ?        00:00:00 bash
   36 ?        00:00:00 bash
   48 ?        00:00:00 ps
```

发现有2个bash进程，分别是启动容器运行的bash和exec启动的bash，容器启动是bash的pid在容器中是1。


到Linux主机的cgroup目录下，利用find命令，找到容器cpu相关的目录：


```
[/sys/fs/cgroup/cpuset/system.slice/docker-30a1da397ec45034cca1582c6a2be9095693a02196a5e40a531ffd4abda1f33d.scope]$ cat cgroup.procs
3022
3381
[/sys/fs/cgroup/cpuset/system.slice/docker-30a1da397ec45034cca1582c6a2be9095693a02196a5e40a531ffd4abda1f33d.scope]$ cat tasks
3022
3381
[/sys/fs/cgroup/cpuset/system.slice/docker-30a1da397ec45034cca1582c6a2be9095693a02196a5e40a531ffd4abda1f33d.scope]$ ps -ef | grep /bin/bash
root      3022  3003  0 13:04 pts/1    00:00:00 /bin/bash
centos    3353  2196  0 13:10 pts/2    00:00:00 /usr/bin/docker-current exec -it 30a1da397ec4 /bin/bash
root      3381  3362  0 13:10 pts/3    00:00:00 /bin/bash
centos    3409  2337  0 13:10 pts/0    00:00:00 grep --color=auto --exclude-dir=.bzr --exclude-dir=CVS --exclude-dir=.git --exclude-dir=.hg --exclude-dir=.svn --exclude-dir=.idea --exclude-dir=.tox /bin/bash
```

通过`cgroup.procs`或者`tasks`文件（记录加入到cgroup的容器进程号），也能看到2个进程，然后通过ps把这2个进程找出来，就是2个`/bin/bash`，这就是容器进程，在Linux主机中的进程号，是通过PID Namespace实现的。


## 资料

未看资料：
超强介绍：https://juejin.im/entry/6844903622698860551

googlesource上介绍了cgroup中各subsystem的各指标的含义：https://kernel.googlesource.com/pub/scm/linux/kernel/git/glommer/memcg/+/cpu_stat/Documentation/cgroups