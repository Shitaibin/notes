[*Docker目录*](https://github.com/Shitaibin/notes/tree/master/docker#%E7%9B%AE%E5%BD%95)

## Cgroup

Cgroup 是 Control Group 的缩写，提供对一组进程，及未来子进程的资源限制、控制、统计能力，包括CPU、内存、磁盘、网络。

Cgroup 包含3个组件：
- cgroup ：一组进程，可以加上subsystem
- subsystem ：一组资源控制模块，CPU、内存...
- hierarchy ： 把一组cgroup串成树状结构，这样就能实现cgroup的继承。为什么要继承呢？就如同docker镜像的继承，站在前人的基础之上，免去重复的配置

### 查看Docker的Cgroup信息

1. 创建一个容器，限制为内存为128MB

```
[~]$ docker run -itd -m 128m ubuntu:16.04
9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529
```

`9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529`为容器id

2. Linux下容器的cgroup信息在`/sys/fs/cgroup/`目录，每个目录应该都是一类资源信息

```
[~]$ cd /sys/fs/cgroup/
[/sys/fs/cgroup]$ ll
总用量 0
drwxr-xr-x 4 root root  0 8月  26 12:10 blkio
lrwxrwxrwx 1 root root 11 8月  26 12:10 cpu -> cpu,cpuacct
lrwxrwxrwx 1 root root 11 8月  26 12:10 cpuacct -> cpu,cpuacct
drwxr-xr-x 4 root root  0 8月  26 12:10 cpu,cpuacct
drwxr-xr-x 3 root root  0 8月  26 12:10 cpuset
drwxr-xr-x 4 root root  0 8月  26 12:10 devices
drwxr-xr-x 3 root root  0 8月  26 12:10 freezer
drwxr-xr-x 3 root root  0 8月  26 12:10 hugetlb
drwxr-xr-x 4 root root  0 8月  26 12:10 memory
lrwxrwxrwx 1 root root 16 8月  26 12:10 net_cls -> net_cls,net_prio
drwxr-xr-x 3 root root  0 8月  26 12:10 net_cls,net_prio
lrwxrwxrwx 1 root root 16 8月  26 12:10 net_prio -> net_cls,net_prio
drwxr-xr-x 3 root root  0 8月  26 12:10 perf_event
drwxr-xr-x 4 root root  0 8月  26 12:10 pids
drwxr-xr-x 4 root root  0 8月  26 12:10 systemd
```

3. 利用容器id找到与当前容器相关的cgroup信息

```
[/sys/fs/cgroup]$ find . -name "*9fe91cc60ff0a94345e13a9*" -print
./memory/system.slice/var-lib-docker-containers-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529-shm.mount
./memory/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
./freezer/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
./blkio/system.slice/var-lib-docker-containers-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529-shm.mount
./blkio/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
./cpuset/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
./hugetlb/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
./devices/system.slice/var-lib-docker-containers-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529-shm.mount
./devices/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
./perf_event/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
./net_cls,net_prio/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
./pids/system.slice/var-lib-docker-containers-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529-shm.mount
./pids/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
./cpu,cpuacct/system.slice/var-lib-docker-containers-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529-shm.mount
./cpu,cpuacct/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
./systemd/system.slice/var-lib-docker-containers-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529-shm.mount
./systemd/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope
```

4. 容器memory的信息在`memory`目录

通过容器id可以找到容器对于的目录：`docker-xxx`：

```
[/sys/fs/cgroup/memory/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope]$ ls
cgroup.clone_children           memory.kmem.slabinfo                memory.memsw.failcnt             memory.soft_limit_in_bytes
cgroup.event_control            memory.kmem.tcp.failcnt             memory.memsw.limit_in_bytes      memory.stat
cgroup.procs                    memory.kmem.tcp.limit_in_bytes      memory.memsw.max_usage_in_bytes  memory.swappiness
memory.failcnt                  memory.kmem.tcp.max_usage_in_bytes  memory.memsw.usage_in_bytes      memory.usage_in_bytes
memory.force_empty              memory.kmem.tcp.usage_in_bytes      memory.move_charge_at_immigrate  memory.use_hierarchy
memory.kmem.failcnt             memory.kmem.usage_in_bytes          memory.numa_stat                 notify_on_release
memory.kmem.limit_in_bytes      memory.limit_in_bytes               memory.oom_control               tasks
memory.kmem.max_usage_in_bytes  memory.max_usage_in_bytes           memory.pressure_level
```

查看内存限制，刚好是设置的128MB：

```
[/sys/fs/cgroup/memory/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope]$ cat memory.limit_in_bytes
134217728
[/sys/fs/cgroup/memory/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope]$ python
Python 2.7.5 (default, Apr 11 2018, 07:36:10)
[GCC 4.8.5 20150623 (Red Hat 4.8.5-28)] on linux2
Type "help", "copyright", "credits" or "license" for more information.
>>> 134217728 / 1024 /1024
128
>>>
```

查看容器内存统计信息：

```
[/sys/fs/cgroup/memory/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope]$ cat memory.stat
cache 0
rss 479232
rss_huge 0
mapped_file 0
swap 0
pgpgin 554
pgpgout 437
pgfault 1862
pgmajfault 0
inactive_anon 0
active_anon 479232
inactive_file 0
active_file 0
unevictable 0
hierarchical_memory_limit 134217728
hierarchical_memsw_limit 268435456
total_cache 0
total_rss 479232
total_rss_huge 0
total_mapped_file 0
total_swap 0
total_pgpgin 554
total_pgpgout 437
total_pgfault 1862
total_pgmajfault 0
total_inactive_anon 0
total_active_anon 479232
total_inactive_file 0
total_active_file 0
total_unevictable 0
[/sys/fs/cgroup/memory/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope]$
[/sys/fs/cgroup/memory/system.slice/docker-9fe91cc60ff0a94345e13a90ad08d2683bf14e2aa37c4d277f384e3bfff17529.scope]$ docker stats
CONTAINER           CPU %               MEM USAGE / LIMIT   MEM %               NET I/O             BLOCK I/O           PIDS
9fe91cc60ff0        0.00%               468 KiB / 128 MiB   0.36%               1.3 kB / 648 B      0 B / 0 B           1
```

可以看到内存统计信息和`docker stats`的信息是匹配的，因为这就是`docker stats`数据的来源。

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

**Cgroup的本质就是通过树形文件系统实现了cgroup hierarchy，cgroup下的目录，都是上层目录的子节点，即子cgroup，会继承上层目录的限制，代表容器的目录，就是针对当前容器的子cgroup，它限定了容器的某一类资源，正是由于这是一颗资源限制树，按不同的资源类别划分，并进行继承，所以一个容器的多种资源限制，分布在多个目录中**。


