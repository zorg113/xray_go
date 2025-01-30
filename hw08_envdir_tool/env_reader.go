package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"
)

const (
	disableEnvNameChars = "="
	trimChars           = " \t"
	MaxEnvVarSize       = 512
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(Environment)
	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, errWalk error) error {
			if errWalk != nil {
				return errWalk
			}
			if !info.Mode().IsRegular() {
				return nil
			}
			if info.Size() > MaxEnvVarSize {
				return nil
			}
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			var val EnvValue
			if info.Size() == 0 {
				val.NeedRemove = true
			} else {
				reader := bufio.NewReader(file)
				line, _, err := reader.ReadLine()
				if err != nil {
					return err
				}
				line = bytes.ReplaceAll(line, []byte{0x00}, []byte{'\n'})
				val.NeedRemove = false
				val.Value = strings.TrimRight(string(line), trimChars)
			}
			envs[filepath.Base(path)] = val
			return nil
		})
	if err != nil {
		return nil, err
	}
	return envs, nil
}
