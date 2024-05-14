package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func main() {
	err := some()
	fmt.Printf("%+v\n", err)
	// fmt.Printf("%+v\n", errors.Cause(err))
}

func some() error {
	if err := a(); err != nil {
		return errors.WithMessage(err, "some error")
	}
	return nil
}

func a() error {
	err := errors.New("test error")
	// return errors.Wrapf(err, "a failed, error=%v", err)
	return errors.WithMessage(err, "a error")
}
