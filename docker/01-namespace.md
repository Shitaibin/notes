[*Docker目录*](https://github.com/Shitaibin/notes/tree/master/docker#%E7%9B%AE%E5%BD%95)

-------------

## 目录

- [目录](#目录)
- [Namespace](#namespace)
	- [实践](#实践)
  
## Namespace

Namespace帮助进程隔离出空间，有以下几个类别的空间：
- PID ： 让每个进程在Namespace中，以pid 1开始，并且与其他Namespace的进程隔离
- UID ：让每个进程拥有在Namespace中有Root权限
- Mount ： 让每个进程拥有虚拟的磁盘挂载
- Network ： 让每个进程拥有虚拟的网络
- UTS ： 让每个Namespace有自己的hostname
- IPC ： 让每个进程的消息队列隔离



### 实践

`echo $$` : 输出当前进程号。

```go
package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
```

- `pstree -pl` ： 列出当前的进程调用树。
- `readlink /proc/PID/ns/uts` : 可以查看某个进程的uts
- `ipcs -q` ： 可以查看系统当前的消息队列
- `ipcmk -Q` ： 可以创建消息
- `mount -t proc proc /proc` : 用于挂载进程空间 , 退出也执行相同命令
- `id` : 命令用来查看当前系统的用户信息