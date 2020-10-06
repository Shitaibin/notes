FROM centos:7

# 设置环境变量指令放前面
ENV PATH /usr/local/bin:$PATH

# 替换成163的源，阿里云的包GPG密钥问题
RUN mv /etc/yum.repos.d/CentOS-Base.repo /etc/yum.repos.d/CentOS-Base.repo.backup
RUN curl -o /etc/yum.repos.d/CentOS-Base.repo http://mirrors.163.com/.help/CentOS7-Base-163.repo
RUN yum makecache

# 安装软件指令放前面
RUN yum install -y make
# 把业务软件的配置,版本等经常变动的步骤放最后



