package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func createChroot() string {
	dir, err := os.MkdirTemp("", "chroot")
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("mkdir -p %s/usr/local/bin && cp /usr/local/bin/docker-explorer %s/usr/local/bin/docker-explorer", dir, dir))
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	return dir
}

// Usage: your_docker.sh run <image> <command> <arg1> <arg2> ...
func main() {
	switch command := os.Args[1]; command {
	case "run":
		dir := createChroot()
		defer os.RemoveAll(dir)
		cmd := exec.Command("chroot", append([]string{dir}, os.Args[3:]...)...)
		cmd.Dir = dir
		cmd.Stdout = os.Stdout
		// This is the key addition for PID namespacing
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Cloneflags: syscall.CLONE_NEWPID, // Create a new PID namespace
		}

		/*
			cmd.SysProcAttr is a field of the exec.Cmd struct that allows you to
			control the system call that's used to create the new process.

			syscall.SysProcAttr is a struct that holds attributes that will be
			applied to the new process.  One of these attributes is Cloneflags,
			which allows you to specify flags to be passed to the clone() system call.
			For basic process creation, fork() is sufficient. But when you need more control
			over the relationship between processes, especially when dealing with containers and isolation, clone() is essential.
		*/
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("%v\n", err)
		}
		os.Exit(cmd.ProcessState.ExitCode())
	default:
		panic("mydocker: '" + command + "' is not a mydocker command.")
	}
}
