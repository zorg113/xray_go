package logger

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func New(level string, path string) (*Logger, error) {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, fmt.Errorf("cannot parse level: %w", err)
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("cannot open file: %v", err)
	}
	logger := logrus.New()
	logger.SetLevel(lvl)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(file)
	return &Logger{logger: logger}, nil

}

func (l Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l Logger) Debug(msg string) {
	l.logger.Debug(msg)
}
