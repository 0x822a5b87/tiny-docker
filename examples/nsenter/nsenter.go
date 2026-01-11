package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		panic("Usage: nsenter <pid>")
	}
	pid := args[1]
	fmt.Printf("pid : {%s}\n", pid)

	command := "nsenter"
	path, err := exec.LookPath(command)
	if err != nil {
		panic(err)
	}
	if err = syscall.Exec(path, []string{command, "-t", pid, "-a", "/bin/bash"}, syscall.Environ()); err != nil {
		panic(err)
	}
}
