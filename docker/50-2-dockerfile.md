[*Docker目录*](https://github.com/Shitaibin/notes/tree/master/docker#%E7%9B%AE%E5%BD%95)

Dockerfile常用
-------------

## 目录

- [目录](#目录)

## 书写规则

Dockerfile会构建缓存，要把不经常变的部分放到前面，而易变的放到后面，这样更改Dockerfile后可以有效利用之前的缓存。

比如ENV放到最前面，然后是安装一些软件包，之后才是业务常变动部分。

## 替换国内镜像源

```Dockerfile
FROM centos:7

# 设置环境变量指令放前面
ENV PATH /usr/local/bin:$PATH

# 替换成阿里云的源
RUN mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
RUN curl -o /etc/yum.repos.d/CentOS-Base.repo https://mirrors.aliyun.com/repo/Centos-6.repo
RUN yum makecache

# 安装软件指令放前面
RUN yum install -y make
# 把业务软件的配置,版本等经常变动的步骤放最后
```