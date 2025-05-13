package main

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go
// https://balun-team.yonote.ru/share/1bae2018-93a9-4710-b536-bdd7087ad8f9/doc/domashnee-zadanie-8-Fx7uvPPgtS

type MultiError struct {
	errs []error
}

func (e *MultiError) Error() string {
	if e == nil || len(e.errs) == 0 {
		return ""
	}
	var builder strings.Builder

	builder.WriteString(strconv.Itoa(len(e.errs)) + " errors occurred:\n")
	for _, err := range e.errs {
		builder.WriteString("\t* " + err.Error())
	}
	builder.WriteString("\n")

	return builder.String()
}

func Append(err error, errs ...error) *MultiError {
	if err == nil {
		return &MultiError{errs: errs}
	}

	var multiError *MultiError
	if !errors.As(err, &multiError) {
		return &MultiError{
			errs: errs,
		}
	}

	multiError.errs = append(multiError.errs, errs...)

	return multiError
}

func TestMultiError(t *testing.T) {
	var err error
	err = Append(err, errors.New("error 1"))
	err = Append(err, errors.New("error 2"))

	expectedMessage := "2 errors occurred:\n\t* error 1\t* error 2\n"
	assert.EqualError(t, err, expectedMessage)
}
