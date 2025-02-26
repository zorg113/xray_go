package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

var ErrWriteRune = errors.New("error write rune")

type stringsMock struct{}

func (sm stringsMock) WriteRune(_ rune) (int, error) {
	return 0, ErrWriteRune
}

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "d\n5abc", expected: "d\n\n\n\n\nabc"},
		{input: "dㄆ4abc", expected: "dㄆㄆㄆㄆabc"},
		// uncomment if task with asterisk completed
		// {input: `qwe\4\5`, expected: `qwe45`},
		// {input: `qwe\45`, expected: `qwe44444`},
		// {input: `qwe\\5`, expected: `qwe\\\\\`},
		// {input: `qwe\\\3`, expected: `qwe\3`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackInvalidString(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidString), "actual error %q", err)
		})
	}
}

func TestStateMashineErrorAtoi(t *testing.T) {
	builder := strings.Builder{}
	state := newSateMashine(&builder)
	state.state = Symbol
	r := rune('৩')
	t.Run(string(r), func(t *testing.T) {
		err := state.checkDigit(r)
		require.Truef(t, errors.Is(err, strconv.ErrSyntax), "actual error %q", err)
	})
}

func TestStateMashineErrorWriteRune(t *testing.T) {
	builder := stringsMock{}
	state := newSateMashine(&builder)
	state.state = Symbol
	r := rune('9')
	t.Run(string(r), func(t *testing.T) {
		err := state.checkDigit(r)
		require.Truef(t, errors.Is(err, ErrWriteRune), "actual checkDigit WriteRune error %q", err)
		err = state.checkSymbol('A')
		require.Truef(t, errors.Is(err, ErrWriteRune), "actual checkSymbol WriteRune error %q", err)
		err = state.end()
		require.Truef(t, errors.Is(err, ErrWriteRune), "actual end WriteRune error %q", err)
	})
}
