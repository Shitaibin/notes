package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const CgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
	// if os.Args[0] == "/proc/self/exe" {
	{
		fmt.Printf("Current pid: %d", syscall.Getpid())

		// 施加内存压力，200MB
		cmd := exec.Command("sh", "-c", "stress --vm-bytes 200m --vm-keep -m 1")
		cmd.SysProcAttr = &syscall.SysProcAttr{
			// Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
		}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal(err)
			os.Exit(-1)
		}
	}

	cmd := exec.Command("/bin/bash")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	cmdPid := cmd.Process.Pid
	fmt.Printf("%d\n", cmdPid)

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
