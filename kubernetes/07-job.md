[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----

## 目录

- [Job](#Job)
- [Job示例](#Job示例)
- [Job原理](#Job原理)
- [CronJob](#CronJob)


## Job

Job与Replicaset类似，也是Pod的Ownereference。

Deployment 和 Replicaset 部署都是永远不会停止的应用/任务，比如服务端，如果Pod内进程停止，Pod会被重新拉起来。

但其实也有这种需求：运行某个**可以结束的任务**，任务进程正常结束后，Pod也就结束了，不需要重启，如果异常结束，可以像Replicaset，重启拉起Pod。

k8s的Job和Linux的Job类似，linux中的命令就可以提交到OS的一个job，比如`echo "hello"`，Linux会创建一个进程执行`echo "hello"`，`cat pod.log | grep "ERROR"`也是1个Job，这里面就涉及到多个进程，把cat进程的结果，交给grep进程。

```sh
# top & 就是创建一个后台运行的job
➜  kubernetes git:(master) ✗ top &
[1] 65815
[1]  + 65815 suspended (tty output)  top
# 通过jobs命令查看当前job
➜  kubernetes git:(master) ✗ jobs
[1]  + suspended (tty output)  top
```

Job可以：
- 顺序运行多个Pod，前一个完成，才运行下一个
- 并行运行Pod，限制并行数量
- 拉起异常退出的Pod

[↑top](#目录)

## Job示例

格式依然是apiVersion、kind、metadata、spec，4个部分。

这是官方计算圆周率的job。

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
```

操作实例：

```sh
# 创建Job
➜  kubernetes git:(master) ✗ kubectl apply -f https://k8s.io/examples/controllers/job.yaml
job.batch/pi created
# 查询Job
➜  kubernetes git:(master) ✗ kubectl get jobs
NAME   COMPLETIONS   DURATION   AGE
pi     0/1           86s        86s
# 查询pods
➜  kubernetes git:(master) ✗ kubectl get pods
NAME                                READY   STATUS              RESTARTS   AGE
pi-cx568                            0/1     ContainerCreating   0          8m7s

# 删除
➜  kubernetes git:(master) ✗ kubectl delete job pi
job.batch "pi" deleted
```

[↑top](#目录)


## Job原理

Job是有Job Controller实现的。

Job Controller向API Server注册，收到事件时，检查活动的Pod数量是否大于并发约束，大于则减少Pod，小于则创建Pod，结果再通知到API Server。

![](http://img.lessisbetter.site/k8s-job.png)


## CronJob

Linux有定时任务，k8s也有定时Job，它就是CronJob。

CronJob的定时格式与Linux `Cron`类似，比如：

```
*/1 * * * * *
```

更多资料参考[这里](https://kubernetes.io/docs/concepts/workloads/controllers/cron-jobs/)。


[↑top](#目录)
