

1. 配置文件增加log的配置：

```ini
[log]
level = 6
filename = project.log
maxsize = 1024
```

更多配置项见beego项目的`fileLogWriter`定义。

2. 启动时，设置log。

```go
// main.go
func init() {
    initLog()
    // other things ...
}

func initLog() {
	logConf := make(map[string]interface{})
	logConf["filename"] = beego.AppConfig.String("log::filename")
	level, _ := beego.AppConfig.Int("log::level")
	logConf["level"] = level
	size, _ := beego.AppConfig.Int("log::maxsize")
	logConf["maxsize"] = size

	confByte, err := json.Marshal(logConf)
	if err != nil {
		fmt.Println("marshal failed,err:", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(confByte))
	logs.SetLogFuncCall(true)
}
```