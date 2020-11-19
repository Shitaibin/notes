创建和初始化kubebuilder项目。




```
mkdir guestbook && cd guestbook
go mod init github.com/shitaibin/guestbook
kubebuilder init --domain github.com/shitaibin/guestbook
```

创建输出：

1. 生成脚手架代码
2. 拉依赖
3. 通过make构建，生成的可执行程序成为`manager`。

```
[~/kubebuilder/guestbook]$ kubebuilder init --domain  github.com/shitaibin/guestbook
Writing scaffold for you to edit...
Get controller runtime:
$ go get sigs.k8s.io/controller-runtime@v0.5.0
go: downloading sigs.k8s.io/controller-runtime v0.5.0
go: downloading k8s.io/client-go v0.17.2
go: downloading k8s.io/apimachinery v0.17.2
go: downloading github.com/imdario/mergo v0.3.6
go: downloading k8s.io/api v0.17.2
go: downloading github.com/golang/groupcache v0.0.0-20180513044358-24b0969c4cb7
go: downloading github.com/prometheus/client_model v0.0.0-20190129233127-fd36f4220a90
go: downloading golang.org/x/crypto v0.0.0-20190820162420-60c769a6c586
go: downloading golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456
go: downloading github.com/prometheus/procfs v0.0.2
go: downloading k8s.io/apiextensions-apiserver v0.17.2
Update go.mod:
$ go mod tidy
go: downloading github.com/onsi/gomega v1.8.1
go: downloading gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127
go: downloading go.uber.org/atomic v1.3.2
go: downloading github.com/fsnotify/fsnotify v1.4.7
go: downloading golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
Running make:
$ make
/home/ubuntu/gopath/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go build -o bin/manager main.go
Next: define a resource with:
$ kubebuilder create api
```

生成API脚手架代码：

```
[~/kubebuilder/guestbook]$ kubebuilder create api --group webapp --version v1 --kind Guestbook
Create Resource [y/n]
y
Create Controller [y/n]
y
Writing scaffold for you to edit...
api/v1/guestbook_types.go
controllers/guestbook_controller.go
Running make:
$ make
/home/ubuntu/gopath/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
go build -o bin/manager main.go
```

kubebuilder生成的项目情况：

```
$ tree guestbook -L 1
.
├── Dockerfile
├── Makefile // 包含构建、生成镜像等
├── PROJECT // 项目信息
├── api // 待实现的api，是create api命令生成的
├── bin // manager生成目录
├── config // 部署到k8s集群所需要的配置，包括crd、rbac等
├── controllers // controller的实现，主要是reconciler
├── go.mod
├── go.sum
├── hack
└── main.go // 创建mgr和reconciler，并绑定等
```



开发operator做哪些事？

- 使用kubebuilder创建项目
- 填写api中的spec（期望样子）和status（资源的状态）
- 填写controller


