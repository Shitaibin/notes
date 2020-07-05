

让Beego使用Mysql存储数据，可以从Mysql查询。

### MySQL先存储测试数据

```sql
create database imooc;

CREATE TABLE `user` (
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `name` varchar(128)    NOT NULL DEFAULT '',
    `gender` tinyint(4) NOT NULL DEFAULT '0',
    `age` int(11) NOT NULL DEFAULT '0',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

insert into user (name,gender,age) values('zhangsan',1,10);
insert into user (name,gender,age) values('lisi',0,11);
insert into user (name,gender,age) values('wangwu',1,12);
```

### 利用bee建立视图

```
bee generate scaffold user -fields="id:int64,name:string,gender:int,age:int" -driver=mysql -conn="root:@tcp(127.0.0.1:3306)/immoc"
```

- scaffold: 脚手架
- fields：字段
- driver：数据库驱动
- conn：连接数据库的字符串，格式："账号:密码:@tcp(ip:port)/数据库名称"

### 增加User的路由

```go
// router.go
beego.Include(&controllers.UserController{})
```

### 增加数据库连接

在main函数的`beego.Run()`前增加，或增加init函数内容为：

```go
// main.go
if err := orm.RegisterDataBase("default", "mysql", "root:rootroot@tcp(127.0.0.1:3306)/imooc?charset=utf8"); err != nil {
		logs.Error(err.Error())
	}
```

遇到错误：

```
[E] [proc.go:203]  register db Ping `default`, Error 1049: Unknown database 'immoc'
```

MySQL升级8.0以上版本后，在用第三方库github.com/Go-SQL-Driver/MySQL打开数据库时会报错`this authentication plugin is not supported`，这是因为MySQL8.0版本修改了加密方式，需要密码的加密方式。

可以看到root用户的加密方式是： `caching_sha2_password` 。

```mysql
mysql> select host,user,plugin from mysql.user;
+-----------+------------------+-----------------------+
| host      | user             | plugin                |
+-----------+------------------+-----------------------+
| localhost | mysql.infoschema | caching_sha2_password |
| localhost | mysql.session    | caching_sha2_password |
| localhost | mysql.sys        | caching_sha2_password |
| localhost | root             | caching_sha2_password |
+-----------+------------------+-----------------------+
4 rows in set (0.00 sec)
```

修改密码，然后可以看到加密方式已变更为：`mysql_native_password` 。
```mysql
mysql> alter user root@localhost identified with mysql_native_password by 'rootroot';
Query OK, 0 rows affected (0.01 sec)

mysql> select host,user,plugin from mysql.user;
+-----------+------------------+-----------------------+
| host      | user             | plugin                |
+-----------+------------------+-----------------------+
| localhost | mysql.infoschema | caching_sha2_password |
| localhost | mysql.session    | caching_sha2_password |
| localhost | mysql.sys        | caching_sha2_password |
| localhost | root             | mysql_native_password |
+-----------+------------------+-----------------------+
4 rows in set (0.00 sec)
```

### 验证

访问8080端口就可以获取所有的User信息：

```json
[
  {
    "Id": 1,
    "Name": "zhangsan",
    "Gender": 1,
    "Age": 10
  },
  {
    "Id": 2,
    "Name": "lisi",
    "Gender": 0,
    "Age": 11
  },
  {
    "Id": 3,
    "Name": "wangwu",
    "Gender": 1,
    "Age": 12
  }
]
```