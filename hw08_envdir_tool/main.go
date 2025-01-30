package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("args are not provided: <env_dir> <command ...>")
		return
	}
	dir, command := os.Args[1], os.Args[2:]
	envs, err := ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	exitCode := RunCmd(command, envs)
	os.Exit(exitCode)
}
