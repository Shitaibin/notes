## 架构

### keda-metrics-apiserver / keda adapter

- Watch 所有的ScaleObject，为它们分别建立scaler instance，在API server会访问 adapter获取指标，adapter 使用metrics中的label找到对应的scaler instance，从scaler instance获取观测的指标。

### keda-operator

#### ScaleObject controller

- 管理HPA
  - 生成HPA对象
    - MetricSpec：通过Scaler获取MetricSpec，添加KEDA的自定义label: scaleObjectName，以便HPA来MetricsAdapter 请求指标时，MetricsAdapter知道是请求哪个对象的指标，避免不同scaleObject的相同指标名称混淆
    - behavior：如果scaleObject未指定，则为空
    - ownerRef：指向scaleObject
    - label：HPA名称超过63会被截断，这里记录HPA非截断的全称、版本等信息
  - Create or
  - Update
    - 对比 HPA 的spec和labels，如果不同则更新
- 管理ScaleObject
  - requestScaleLoop
    - 取消已存在了 scaleLoop
    - go startScaleLoop
      - 这是内部Scaler走的逻辑，周期性调用 ScalerActive 检查事件源是否还Active，依据Active状态，记录数据或进行Scale扩缩
        - 获取workload当前currentReplica
        - active
          - currentReplica=0，需要从0扩上来到minReplica（无则为1）
          - 记录 LastActiveTime
        - inactive
          - currentReplica>0，需要缩到idleReplica或minReplica（无则为0），前提时距LastActiveTime已经超过 cooldownPeriod，方式是直接修改Scale子资源的Replica
          - ScaleObject Active Condition设置为false
    - go startPushScalers
      - 如果存在ExternalPushScaler，则执行ExternalPushScaler.Run，通过gRPC和外部的Scaler进行通信，获取外部外部Scaler所关注的事件源是否还Active，然后进行处理。处理同startScaleLoop。
      - 读active channel后，根据active和scaleObject处理：ScaleExecutor 

#### Scaler

- GetMetrics: HPA 要使用的指标当前的数值
- GetMetricSpecForScaling： 生成HPA使用的metric spec
- IsActive：inactive到active时需要执行从0到1的弹性，相反执行从1到0的弹性

### 疑问

1. 缩到0时，scaleObj直接修改scale subresource的replica，HPA没有删除，HPA怎么没有把replica调成hpa的targetReplica？
   1. 猜测：
      1. hpa的target 同 idle或min？概率小。
      2. ？

### Deploy

```
NAME                                      READY   STATUS    RESTARTS   AGE
keda-metrics-apiserver-77b775f485-bbjlz   1/1     Running   0          4m43s
keda-operator-96f9b5767-zfv2x             1/1     Running   0          4m43s
```

```
NAME                                    CREATED AT
clustertriggerauthentications.keda.sh   2022-02-19T14:04:08Z
scaledjobs.keda.sh                      2022-02-19T14:04:08Z
scaledobjects.keda.sh                   2022-02-19T14:04:08Z
triggerauthentications.keda.sh          2022-02-19T14:04:08Z
```

### 值得学习的地方？

1. 