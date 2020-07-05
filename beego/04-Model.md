04-Model.md


## beego orm 

model中，结构体的定义如果没有`Id`字段，会报如下日志，并超时。因为beego orm默认使用Id作为主键。

```
<orm.RegisterModel> `guess/models.User` needs a primary key field, default is to use 'id' if not set
```

### 数据库连接

在main.go中增加数据库连接。

```go
func init() {
	if err := orm.RegisterDriver("mysql", orm.DRMySQL); err != nil {
		logs.Error(err.Error())
	}
	if err :=orm.RegisterDataBase("default", "mysql", "root:rootroot@/guess?charset=utf8"); err != nil {
Error**(err.Error())
	}
}
```

### Beego 参数

可以从Controller中提取。




```go
// 获取id的参数
// http://domain?id=1
id, err := c.GetInt("id")
```