package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger LoggerConf
	DB     DBConf
	Server HTTPConf
}

type LoggerConf struct {
	Level string
	Path  string
}

type HTTPConf struct {
	Host string
	Port string
}

type DBConf struct {
	User     string
	Password string
	Host     string
	Port     uint64
	Name     string
}

func NewConfig(path string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(path, &config); err != nil {
		fmt.Println(err)
		return Config{}, err
	}
	return config, nil
}
