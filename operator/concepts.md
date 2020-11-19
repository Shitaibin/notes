编写operator时的一些概念：

1. Scheme：定义了序列化和反序列化API对象的方法，即它提供了 Kinds 和相应的 Go 类型之间的映射。

2. [*manager*](https://pkg.go.dev/sigs.k8s.io/controller-runtime/pkg/manager#Manager)，它记录着我们所有控制器的运行情况，以及设置共享缓存和API服务器的客户端（注意，我们把我们的 Scheme 的信息告诉了 manager）
3. 开发operator本质是完成controller，其中一个重要的接口是`Reconciler`，即调和，它是controller的核心组成部分，controller接收到API请求后，利用reconciler把集群状态调和到期望的状态。



resources 总是用小写，按照惯例是 Kind 的小写形式。比如pods是资源，Pod是Kind（类型），类型和资源通常一一对应。

- GVK = Group Version Kind ，指某个api group下特定版本的类型。
- GVR = Group Version Resources，指某个api group下特定版本的资源。

所以每个GVK和GVR通常也是一一对应的。

GVK如何体现在Go代码中呢，会有对应的Go type，比如struct，举个例子：

- Go Type：`"tutorial.kubebuilder.io/api/v1".CronJob{}`，包`tutorial.kubebuilder.io/api/v1`下有个结构体`CronJob`

- GVK:  `batch.tutorial.kubebuilder.io/v1`，api group是`batch.tutorial.kubebuilder.io`，版本是`v1`



TypeMeta ：描述了API版本和种类，
ObjectMeta ：包含名称,名称空间和标签等

开发规范：
1. 字段可以使用大多数的基本类型。数字是个例外：出于 API 兼容性的目的，我们只允许三种数字类型。对于整数，需要使用 int32 和 int64 类型；对于小数，使用 resource.Quantity 类型（用10进制或2进制表示浮点数，比如2k=2000,  2m=0.002是10进制表示，2Ki=2048是2进制表示, 2.5的写法是2500m）。
2. 还有一个我们使用的特殊类型：metav1.Time。 它有一个稳定的、可移植的序列化格式的功能，其他与 time.Time 相同