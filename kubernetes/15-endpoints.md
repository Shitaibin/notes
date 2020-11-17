[*k8s目录*](https://github.com/Shitaibin/notes/tree/master/kubernetes#%E7%9B%AE%E5%BD%95)

----


## 目录

- [目录](#目录)
- [Endpoints](#endpoints)
- [样例](#样例)
  - [代表Pod的Endpoints](#代表pod的endpoints)
  - [外部服务的Endpoint](#外部服务的endpoint)

## Endpoints

Endpoints（复数）是kubernetes中的一个二级概念，使用者通常是无需直接去操作Endpoints，直到你需要掌握Service等的原理，或者去定位问题。

Endpoint（单数）是网络端点的含义，格式是`ip:port`，在k8s里有2种网络端点：
1. k8s内，一个pod对应1个endpoint
2. k8s外的服务，比如mysql集群、web集群等，每一个可访问的`ip:port`都是一个endpoint

Endpoints是一组endpoint，它所形成的是`service`名字到一组`ip:port`的映射，本质是为名字服务服务的。

Endpoints的**作用**只有1个，service把流量路由到Endpoints中的Endpoint，其实service不知道这个Endpoint对应的是Pod还是集群外的服务，这样就把Service和具体的服务端进行了解耦。

Endpoints的名称与Service的名称相同。

代表Pods列表的Endpoints是k8s自动创建的，k8s的service设置了label selector，通过label selector构建Endpoints，在Pod变动后，Endpoints也会随着更新。

![endpoints for pods](http://img.lessisbetter.site/2020-11-endpoints-for-pod.png)

代表k8s外服务端的Endpoints需要**手动创建**，名称并且与service的名称相同。

![endpoints for outside server](http://img.lessisbetter.site/2020-11-endpoints-for-outside-server.png)

## 样例

### 代表Pod的Endpoints

部署一个[nginx-svc.yaml](examples/04-service/endpoints/nginx-svc.yaml)服务，然后查看服务、Endpoints。

```
$ kubectl apply -f nginx-svc.yaml
deployment.apps/nginx-deployment created
service/ingress-nginx created
$
$ kubectl get endpoints
NAME            ENDPOINTS                                   AGE
ingress-nginx   10.32.0.15:80,10.32.0.16:80,10.32.0.17:80   7s
$
$ kubectl describe svc ingress-nginx
Name:                     ingress-nginx
...
Selector:                 app=nginx
...
Endpoints:                10.32.0.15:80,10.32.0.16:80,10.32.0.17:80
...
$
$ kubectl describe ep ingress-nginx
Name:         ingress-nginx
Namespace:    default
Labels:       <none>
Annotations:  endpoints.kubernetes.io/last-change-trigger-time: 2020-11-06T02:01:42Z
Subsets:
  Addresses:          10.32.0.15,10.32.0.16,10.32.0.17
  NotReadyAddresses:  <none>
  Ports:
    Name  Port  Protocol
    ----  ----  --------
    http  80    TCP

Events:  <none>
```

### 外部服务的Endpoint

在k8s外部署一个nginx，可访问的ip和端口为：`10.10.50.53:80`，为此建立Endpoints和Service，k8s内部即可访问该外部的nginx。

```
$ kubectl apply -f outserver-svc.yaml
service/external-web created
endpoints/external-web created
$
$ kubectl describe svc external-web
Name:              external-web
...
Type:              ClusterIP
IP:                10.111.43.130
Port:              nginx  80/TCP
TargetPort:        80/TCP
Endpoints:         10.10.50.53:80
...

$ kubectl describe ep external-web
Name:         external-web
Namespace:    default
Labels:       <none>
Annotations:  <none>
Subsets:
  Addresses:          10.10.50.53
  NotReadyAddresses:  <none>
  Ports:
    Name   Port  Protocol
    ----   ----  --------
    nginx  80    TCP

Events:  <none>
```
