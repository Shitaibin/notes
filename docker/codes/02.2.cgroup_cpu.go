package main

// CPU子系统测试

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

func main() {
	if os.Args[0] == "/proc/self/exe" {
		fmt.Println("---------- 2 ------------")
		fmt.Printf("Current pid: %d\n", syscall.Getpid())

		// 开启3个goroutine，占用300%的CPU
		for i := 0; i < 3; i++ {
			go work(i)
		}

		select {}
	}

	fmt.Println("---------- 1 ------------")
	// 创建子cpu cgroup和cpuset group
	cpuGroup := path.Join("/sys/fs/cgroup/cpu", "test_cpu_limit")
	cpusetGroup := path.Join("/sys/fs/cgroup/cpuset", "test_cpuset_limit")
	defer func() {
		fmt.Println("Exit main")
		if err := os.RemoveAll(cpuGroup); err != nil {
			fmt.Printf("Remove cpuGroup path error: %v\n", err)
		}
		if err := os.RemoveAll(cpusetGroup); err != nil {
			fmt.Printf("Remove cpusetGroup path error: %v\n", err)
		}
	}()

	if err := os.Mkdir(cpuGroup, 0755); err != nil {
		fmt.Printf("Mkdir cpuGroup error: %v\n", err)
	}
	if err := os.Mkdir(cpusetGroup, 0755); err != nil {
		fmt.Printf("Mkdir cpuGroup error: %v\n", err)
	}

	fmt.Printf("Test type: ")
	if len(os.Args) <= 1 || os.Args[1] == "nolimit" {
		// 无限制
		fmt.Println("No limit")
	} else if os.Args[1] == "cpu" {
		fmt.Println("Cpu limit")
		// 限制cpu使用时间，使用cfs类型
		if err := ioutil.WriteFile(path.Join(cpuGroup, "cpu.cfs_quota_us"), []byte("500"), 0755); err != nil {
			fmt.Printf("Write cpu.cfs_quota_us error: %v\n", err)
			return
		}
	} else if os.Args[1] == "cpuset" {
		fmt.Println("Cpuset limit")
		// 限制使用的cpu核
		if err := ioutil.WriteFile(path.Join(cpusetGroup, "cpuset.cpus"), []byte("1,3"), 0644); err != nil {
			fmt.Printf("Write cpuset.cpus error: %v\n", err)
			return
		}
	} else {
		fmt.Printf("Invalid parameter: %s\n", os.Args[1])
		return
	}

	// 启动子进程
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Printf("/proc/self/exe start error: %v\n", err)
		return
	}

	cmdPid := cmd.Process.Pid
	fmt.Printf("cmdPid: %d\n", cmdPid)
	cmdPidByte := []byte(strconv.Itoa(cmdPid))

	// 进程加入到两个cgroup

	if err := ioutil.WriteFile(path.Join(cpusetGroup, "tasks"), cmdPidByte, 0644); err != nil {
		fmt.Printf("Add task to cpuset cgroup error: %s\n", err)
	}

	if err := ioutil.WriteFile(path.Join(cpuGroup, "tasks"), cmdPidByte, 0644); err != nil {
		fmt.Printf("Add task to cpu cgroup error: %s\n", err)
	}

	cmd.Process.Wait()
}

func work(id int) {
	fmt.Printf("worker %d start\n", id)
	for {
	}
}
