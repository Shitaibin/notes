package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const CgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
		fmt.Println("---------- 2 ------------")
		fmt.Printf("Current pid: %d\n", syscall.Getpid())

		// 创建stress子进程，施加内存压力
		allocMemSize := "101m" // 另外1项测试为99m
		fmt.Printf("allocMemSize: %v\n", allocMemSize)
		stressCmd := fmt.Sprintf("stress --vm-bytes %s --vm-keep -m 1", allocMemSize)
		cmd := exec.Command("sh", "-c", stressCmd)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("stress run error: %v", err)
			os.Exit(-1)
		}
	}

	fmt.Println("---------- 1 ------------")
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWPID,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// 启动子进程
	if err := cmd.Start(); err != nil {
		fmt.Printf("/proc/self/exe start error: %v", err)
		os.Exit(-1)
	}

	cmdPid := cmd.Process.Pid
	fmt.Printf("cmdPid: %d\n", cmdPid)

	// 创建子cgroup
	memoryGroup := path.Join(CgroupMemoryHierarchyMount, "test_memory_limit")
	os.Mkdir(memoryGroup, 0755)
	// 设定内存限制
	ioutil.WriteFile(path.Join(memoryGroup, "memory.limit_in_bytes"),
		[]byte("100m"), 0644)
	// 将进程加入cgroup
	ioutil.WriteFile(path.Join(memoryGroup, "tasks"),
		[]byte(strconv.Itoa(cmdPid)), 0644)

	cmd.Process.Wait()
}
