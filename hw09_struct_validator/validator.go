package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type bound int

const (
	intMax bound = iota
	intMin
)

var (
	ErrTag        = errors.New("validate tag will-formed")
	ErrNotAStruct = errors.New("interface is not a struct")
	ErrLogic      = errors.New("incorrect boundary condition")
)

var (
	ErrMin    = errors.New("value is less than minimum")
	ErrMax    = errors.New("value us greater than maximum")
	ErrIn     = errors.New("no values match")
	ErrLen    = errors.New("string len doesn't fit")
	ErrRegExp = errors.New("string doesn't match regexp")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errsLine := "ValidationErrors: "
	for i := range v {
		errsLine += fmt.Sprintf(" %s", v[i])
	}
	return errsLine
}

func Validate(v interface{}) error {
	var errorsArray ValidationErrors
	obj := reflect.ValueOf(v)
	if obj.Kind() != reflect.Struct {
		return ErrNotAStruct
	}
	t := obj.Type()

	for i := 0; i < t.NumField(); i++ {
		fType := t.Field(i)
		fVal := obj.Field(i)
		validateTag, ok := fType.Tag.Lookup("validate")
		if !ok {
			continue
		}
		conditions := strings.FieldsFunc(validateTag, func(r rune) bool { return r == '|' })
		if len(conditions) == 0 {
			return ErrTag
		}
		for el := range conditions {
			err := ValidateConditions(fVal, conditions[el])
			var errs ValidationErrors
			switch {
			case errors.Is(err, ErrTag):
				return err
			case errors.As(err, &errs):
				errorsArray = append(errorsArray, errs...)
			case err != nil:
				errorsArray = append(errorsArray, ValidationError{Err: err, Field: fType.Name})
			}
		}
	}
	if len(errorsArray) == 0 {
		return nil
	}
	return errorsArray
}

func ValidateConditions(val reflect.Value, condition string) (err error) {
	switch val.Kind() { //nolint:exhaustive
	case reflect.String:
		err = ValidateString(val.String(), condition)
	case reflect.Int:
		err = ValidateInt(int(val.Int()), condition)
	case reflect.Slice:
		for i := 0; i < val.Len(); i++ {
			if err = ValidateConditions(val.Index(i), condition); err != nil {
				break
			}
		}
	case reflect.Struct:
		if condition == "nested" && val.CanInterface() {
			err = Validate(val.Interface())
		}
	default:
		fmt.Printf("not implemented tag: %v\n", val.Kind())
	}
	return
}

func ValidateInt(intVal int, tag string) error {
	args := strings.FieldsFunc(tag, func(r rune) bool { return r == ':' })
	if len(args) != 2 {
		return ErrTag
	}
	switch {
	case args[0] == "min":
		if err := checkLimValue(args[1], intVal, intMin); err != nil {
			return err
		}
	case args[0] == "max":
		if err := checkLimValue(args[1], intVal, intMax); err != nil {
			return err
		}
	case args[0] == "in":
		vals := strings.FieldsFunc(args[1], func(r rune) bool { return r == ',' })
		inValues := make([]int, 0, len(vals))
		for i := range vals {
			v, err := strconv.Atoi(vals[i])
			if err != nil {
				return ErrTag
			}
			inValues = append(inValues, v)
		}
		for i := range inValues {
			if inValues[i] == intVal {
				return nil
			}
		}
		return ErrIn
	default:
		return ErrTag
	}
	return nil
}

func checkLimValue(str string, iVal int, b bound) error {
	i, err := strconv.Atoi(str)
	if err != nil {
		return ErrTag
	}
	switch {
	case b == intMin:
		if i >= iVal {
			return ErrMin
		}
	case b == intMax:
		if i <= iVal {
			return ErrMax
		}
	default:
		return ErrLogic
	}
	return nil
}

func ValidateString(str string, tag string) error {
	args := strings.FieldsFunc(tag, func(r rune) bool { return r == ':' })
	if len(args) != 2 {
		return ErrTag
	}
	switch {
	case args[0] == "len":
		i, err := strconv.Atoi(args[1])
		if err != nil {
			return ErrTag
		}
		if len(str) != i {
			return ErrLen
		}
	case args[0] == "regexp":
		r, err := regexp.Compile(args[1])
		if err != nil {
			return ErrTag
		}
		if !r.MatchString(str) {
			return ErrRegExp
		}
	case args[0] == "in":
		vals := strings.FieldsFunc(args[1], func(r rune) bool { return r == ',' })
		for i := range vals {
			if vals[i] == str {
				return nil
			}
		}
		return ErrIn
	default:
		return ErrTag
	}
	return nil
}
