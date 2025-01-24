package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	disableEnvNameChars = "="
	trimChars           = " \n\t"
)

type Environment map[string]EnvValue

func (e Environment) Strings() []string {
	kenv := make([]string, 0, len(e))
	for key, val := range e {
		kenv = append(kenv, fmt.Sprintf("%v=%v", key, val.Value))
	}
	return kenv
}

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	filesInfo, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read dir error: %w", err)
	}
	envs := make(Environment, len(filesInfo))
	for _, file := range filesInfo {
		if !file.Type().IsRegular() || file.Type().IsDir() {
			continue
		}
		if strings.ContainsAny(file.Name(), disableEnvNameChars) {
			continue
		}
		var fi os.FileInfo
		if fi, err = file.Info(); err != nil {
			return nil, fmt.Errorf("read dir error: %w", err)
		}

		val, err := func(finf os.FileInfo) (string, error) {
			f, err := os.Open(filepath.Join(dir, finf.Name()))
			if err != nil {
				return "", fmt.Errorf("read %v file error: %w", finf.Name(), err)
			}
			defer f.Close()
			r := bufio.NewReader(f)
			var (
				buf        strings.Builder
				bytesValue []byte
				isPrefix   = true
			)
			for isPrefix {
				if bytesValue, isPrefix, err = r.ReadLine(); err != nil && !errors.Is(err, io.EOF) {
					return "", fmt.Errorf("read %v file error: %w", finf.Name(), err)
				}
				buf.Write(bytesValue)
			}
			return buf.String(), nil
		}(fi)
		if err != nil {
			return nil, err
		}
		val = normalizeValue(val)
		if val != "" {
			envs[file.Name()] = EnvValue{val, false}
			fmt.Printf("==> %s %s\n", file.Name(), val)
		}
	}
	return envs, nil
}

func normalizeValue(value string) string {
	if value == "" {
		return value
	}
	value = strings.TrimRight(value, trimChars)
	value = strings.ReplaceAll(value, "\x00", "\n")
	return value
}
