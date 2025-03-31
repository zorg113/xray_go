package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("open log", func(t *testing.T) {
		var err error
		log, err := New("INFO", "/tmp/test")
		require.NoError(t, err)
		log.Info("hello world")
	})
}
