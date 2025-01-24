package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const (
	ExitCodeIOError         = 5
	ExitCodeCommandNotFound = 127
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) < 1 {
		return ExitCodeCommandNotFound
	}
	command, args := cmd[0], cmd[1:]
	execCmd := exec.Command(command, args...)
	execCmd.Stdin, execCmd.Stdout, execCmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)		
		for val := range env.Strings() {
			curr :=  strings.SplitN(val, "=", 2)		
			if pair[0] == curr[0] {
	           
				
				



				
			}
		}
		fmt.Println(pair)
	}	
	if err := execCmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return exitErr.ExitCode()
		}
		return ExitCodeIOError
	}
	return 0
}
