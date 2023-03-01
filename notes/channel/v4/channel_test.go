package main

import (
	"testing"

	"go.uber.org/goleak"
)

func TestChannel(t *testing.T) {
	defer goleak.VerifyNone(t)
	example1()
}
