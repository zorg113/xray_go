package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type States int

const (
	Empty States = iota
	Symbol
)

type writer interface {
	WriteRune(s rune) (int, error)
}

type stateMashine struct {
	state States
	prev  rune   // запоминаемый символ
	out   writer // распакованная сторка (интерфейс)
}

func newSateMashine(outIn writer) *stateMashine {
	return &stateMashine{state: Empty, prev: 0, out: outIn}
}

// поверка входного символа на число и распаковка

func (m *stateMashine) checkDigit(s rune) error {
	n := 0
	err := error(nil)
	if unicode.IsDigit(s) {
		if m.state == Empty {
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
		m.state = Empty
	}
	return nil
}

// проврка на нечисловой символ и сохранение его для
// возможной распаковки

func (m *stateMashine) checkSymbol(s rune) error {
	if !unicode.IsDigit(s) {
		if m.state == Symbol {
			if _, err := m.out.WriteRune(m.prev); err != nil {
				return err
			}
		}
		m.prev = s
		m.state = Symbol
	}
	return nil
}

// сброс последнего состояния и вывод сохраненного последнего символа

func (m *stateMashine) end() error {
	if m.state == Symbol {
		if _, err := m.out.WriteRune(m.prev); err != nil {
			return err
		}
		m.state = Empty
	}
	return nil
}

// распаковка строк

func Unpack(s string) (string, error) {
	builder := strings.Builder{}
	extract := newSateMashine(&builder)
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
	return builder.String(), nil
}
