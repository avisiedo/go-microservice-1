package common

import (
	"errors"
	"fmt"
)

var ErrNotImplemented = errors.New("not implemented")

func checkEmptyFieldName(fieldName string) {
	if fieldName == "" {
		panic(ErrEmpty("fieldName"))
	}
}

// ErrNil create an error for a nil value at the presenter
// component level.
func ErrNil(fieldName string) error {
	checkEmptyFieldName(fieldName)
	return fmt.Errorf("'%s' is nil", fieldName)
}

// ErrEmpty create an error for an empty value at the
// presenter level.
func ErrEmpty(fieldName string) error {
	checkEmptyFieldName(fieldName)
	return fmt.Errorf("'%s' is empty", fieldName)
}
