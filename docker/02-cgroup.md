[*Docker目录*](https://github.com/Shitaibin/notes/tree/master/docker#%E7%9B%AE%E5%BD%95)

资源环境：

```
[~/notes/docker/codes]$ uname -a                                                                                     *[master]
Linux aliyun 4.4.0-117-generic #141-Ubuntu SMP Tue Mar 13 12:01:47 UTC 2018 i686 i686 i686 GNU/Linux
```

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


### 利用Go演示Cgroup内存限制

#### 源码

cgroup的演示[源码](./codes/02.1.cgroup.go) ，关于源码中的`/proc/self/exe`看[补充小知识](#补充小知识)。


源码运行解读：
1. 使用`go run`运行程序，或build后运行程序时，程序的名字是`02.1.cgroup`，所以不满足`os.Args[0] == "/proc/self/exe"`会被跳过。
2. 然后使用`"/proc/self/exe"`新建了子进程，子进程此时叫：`"/proc/self/exe"`
3. 创建cgroup `test_memory_limit`，然后设置内存限制为100MB
4. 把子进程加入到cgroup `test_memory_limit`
5. 等待子进程结束
6. 子进程干了啥呢？子进程其实还是当前程序，只不过它的名字是`"/proc/self/exe"`，符合最初的if语句，之后它会创建stress子进程，然后运行stress，可以修改`allocMemSize`设置stress所要占用的内存

修改源码，stress命令设置内存占用为99m，然后启动测试程序：

#### 不超越内存限制情况

```
[~/workspace/notes/docker/codes]$ go run 02.1.cgroup.go
---------- 1 ------------
cmdPid: 2533
---------- 2 ------------
Current pid: 1
allocMemSize: 99m
stress: info: [6] dispatching hogs: 0 cpu, 0 io, 1 vm, 0 hdd
```

可以看到，子进程`"/proc/self/exe"`运行后取得的pid为 **2533** ，在新的Namespace中，子进程`"/proc/self/exe"`的pid已经变成1，然后利用stress打了99M内存。

使用top查看资源使用情况，stress进程内存RES大约为99M，pid 为 **2539** 。

```
  PID USER      PR  NI    VIRT    RES    SHR S %CPU %MEM     TIME+ COMMAND
 2539 root      20   0  103940 101680    284 R 93.8  9.9   0:06.09 stress
```

```
[/sys/fs/cgroup/memory/test_memory_limit]$ cat memory.limit_in_bytes
104857600
[/sys/fs/cgroup/memory/test_memory_limit]$ # 104857600 刚好为100MB
[/sys/fs/cgroup/memory/test_memory_limit]$ cat memory.usage_in_bytes
2617344
[/sys/fs/cgroup/memory/test_memory_limit]$ cat tasks
2533 <--- /prof/self/exe进程
2534
2535
2536
2537
2538
2539 <--- stress进程
```

tasks下都是在cgroup `test_memory_limit` 中的进程，这些是Host中真实的进程号，通过`pstree -p`查看进程树，看看这些都是哪些进程：

![Cgroup限制内存的进程树](http://img.lessisbetter.site/2020-08-cgroup.png)

进程树佐证了前面的代码执行流程分析大致是对的，只不过这其中还涉及一些创建子进程的具体手段，比如stress是通过sh命令创建出来的。

#### 内存超过限制被Kill情况

内存超过cgroup限制的内存会怎么样？会OOM吗？

如果将内存提高到占用101MB，大于cgroup中内存的限制100M时就会被Kill。

```
[~/notes/docker/codes]$ go run 02.1.cgroup.go                                                                        *[master]
---------- 1 ------------
cmdPid: 21492
---------- 2 ------------
Current pid: 1
allocMemSize: 101m
stress: info: [6] dispatching hogs: 0 cpu, 0 io, 1 vm, 0 hdd
stress: FAIL: [6] (415) <-- worker 7 got signal 9
stress: WARN: [6] (417) now reaping child worker processes
stress: FAIL: [6] (421) kill error: No such process
stress: FAIL: [6] (451) failed run completed in 0s
2020/08/27 17:38:52 exit status 1
```

`stress: FAIL: [6] (415) <-- worker 7 got signal 9` 说明收到了信号9，即SIGKILL 。



### 补充小知识

在演示源码中，使用到`"/proc/self/exe"`，它在Linux是一个特殊的软链接，它指向当前正在运行的程序，比如执行`ll`查看该文件时，它就执行了`/usr/bin/ls`，因为当前的程序是`ls`：

```
[~]$ ll /proc/self/exe
lrwxrwxrwx 1 centos centos 0 8月  27 12:44 /proc/self/exe -> /usr/bin/ls
```

演示代码中的技巧就是通过`"/proc/self/exe"`重新启动一个子进程，只不过进程名称叫`"/proc/self/exe"`而已。如果代码中没有那句if判断，又会执行到创建子进程，最终会导致递归溢出。