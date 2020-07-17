




## 目录

- [目录](#目录)
- [url带入参数](#url带入参数)
	- [方法一](#方法一)
	- [方法二](#方法二)
- [Json携带参数](#json携带参数)
	- [方法一](#方法一-1)
	- [方法二](#方法二-1)


## url带入参数

url带参数适用于GET和POST等操作。

样例url参数：

```
http://localhost/register?name=zw&pwd=123
```

### 方法一

Controller包含了很多获取参数的接口，比如：

```go
c.GetString(key)
c.GetInt(key)
```

利用GetString获取name的值，利用GetInt获取pwd的值。

### 方法二


通过ParseForm解析参数到结构体。

```go
type UserRequest struct {
	Name     string `form:"name"` // 字段要与参数相同，包括大小写，如果不同，使用`from:`说明对应的参数
	Password string `form:"pwd"`
}

var u UserRequest
err := r.ParseForm(&u)
if err != nil {
    r.Ctx.WriteString("read user and password, get error: " + err.Error())
    return
}
```

## Json携带参数

Json携带参数是放在请求Body中的，不适合GET，适合POST等有请求Body的操作。

请求样例：

1. 设置：application/json; charset=utf-8
2. json参数

```json
{
    "name": "dabin",
    "password": "123"
}
```

前提修改配置，把HttpBody带入Context中，配置中增加：

```ini
copyrequestbody = true 
```

### 方法一

解析到结构体里。

```go
func (j *JsonController) Post() {
	// 解析json参数
	var req JsonRequest
	if err := json.Unmarshal(j.Ctx.Input.RequestBody, &req); err != nil {
		j.Ctx.WriteString("json unmarshal error: " + err.Error())
		return
	}

	// 填写json格式响应
	data, _ := json.Marshal(req)
	j.Data["json"] = string(data)
	// 以json格式发送响应
	j.ServeJSON()
}

type JsonRequest struct {
	Name     string
	Password string
}
```

### 方法二

把参数解析到map里。

```go
func (j *JsonController) Post() {
	// 解析json参数
	var req map[string][string]
	if err := json.Unmarshal(j.Ctx.Input.RequestBody, &req); err != nil {
		j.Ctx.WriteString("json unmarshal error: " + err.Error())
		return
	}

	// 填写json格式响应
	data, _ := json.Marshal(req)
	j.Data["json"] = string(data)
	// 以json格式发送响应
	j.ServeJSON()
}
```