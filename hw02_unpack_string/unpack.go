package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type stateMashine struct {
	state int             // варианты состояния и переходов 1 и 0
	prev  rune            // запоминаемый символ
	out   strings.Builder // распакованная сторка
}

func newSateMashine() *stateMashine {
	return &stateMashine{}
}

func (m *stateMashine) checkDigit(s rune) error {
	n := 0
	err := error(nil)
	if unicode.IsDigit(s) {
		if m.state == 0 {
			return ErrInvalidString
		}
		if n, err = strconv.Atoi(string(s)); err != nil {
			return err
		}
		for i := 0; i < n; i++ {
			if _, err = m.out.WriteRune(m.prev); err != nil {
				return err
			}
		}
		m.state = 0
	}
	return nil
}

func (m *stateMashine) checkSymbol(s rune) error {
	if !unicode.IsDigit(s) {
		if m.state == 1 {
			if _, err := m.out.WriteRune(m.prev); err != nil {
				return err
			}
		}
		m.prev = s
		m.state = 1
	}
	return nil
}

func (m *stateMashine) end() error {
	if m.state == 1 {
		if _, err := m.out.WriteRune(m.prev); err != nil {
			return err
		}
	}
	return nil
}

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	extract := newSateMashine()
	for _, val := range s {
		if err := extract.checkDigit(val); err != nil {
			return "", err
		}
		if err := extract.checkSymbol(val); err != nil {
			return "", err
		}
	}
	if err := extract.end(); err != nil {
		return "", err
	}
	return extract.out.String(), nil
}
