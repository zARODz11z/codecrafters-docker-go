package main

import (
	"fmt"
	"os"
	"os/exec"
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
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("%v\n", err)
		}
		os.Exit(cmd.ProcessState.ExitCode())
	default:
		panic("mydocker: '" + command + "' is not a mydocker command.")
	}
}
