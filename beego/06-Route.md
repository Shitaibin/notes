
Beego支持3种路由，样式如下：

```go
beego.Router("/", &controllers.MainController{})
beego.Include(&controllers.CctpInfoController{})
beego.AutoRouter(&controllers.UserController{})
```

### 默认路由

- Router 为默认路由，适合最简单的PUT、GET等操作，非常适合RESTful，如果url中携带参数，比如`user/1`其中1为参数是解析不出来的，需要使用在其中使用正则表达式： `/user/:id`，才能匹配路由和解析参数。

### 注解路由

Include 为注解路由，bee构建时，会从Controller的注释中寻找`@router`，生成该Controller的路由，路由代码自动保存到`router/commentsRouter_controllers.go`中。对于beego生成的脚手架代码，非常适合使用这种方式。



使用注解路由分3步：

第1步：修改配置文件

注解路由只能在dev模式中使用，在`conf/app.conf`中增加：

```ini
runmode = dev
```

同时增加`EnableAdmin`可以在监控`http://localhost:8088/listconf?command=router`中查看所有的路由情况。

```ini
EnableAdmin = true
```


第2步：`router.go`中使用`Include`，增加路由


第3步：Controller的函数增加路由注释

Controller的函数上方注释中增加 `// @router ....`，样例如下


```go
// GetOne ...
// @Title Get One 
// @Description get User by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.User
// @Failure 403 :id is empty
// @router /user/:id [get]
func (c *UserController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetUserById(id)
	if err != nil {
		c.Data["json"] = err.Error()
	} else {
		c.Data["json"] = v
	}
	c.ServeJSON()
}
```

`// @router /user/:id [get]`的含义如下：

- URL `/user/:id` 的GET操作会被路由到`UserController.GetOne`函数。
- `:id`：为路径参数
- [get]：为对应的请求操作类型

第4步：生成注解路由

运行`bee run`，上方的样例在`router/commentsRouter_controllers.go`中生成的注解路由如下：

```go
beego.GlobalControllerRouter["hello-beego/controllers:UserController"] = append(beego.GlobalControllerRouter["hello-beego/controllers:UserController"],
        beego.ControllerComments{
            Method: "GetOne",
            Router: "/user/:id",
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})
```

### 自动路由

- AutoRouter 为自动路由，自动分析Controller拥有的函数，然后把函数转变为路径。比如`UserController.Login()`可以得到`/user/login`这个路由，但无法解析参数。

