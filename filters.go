package main

import (
	"bytes"
)

// filterFlags is a type used in flag.Var to allow repetition of the same
// flag.
type filterFlags []string

// String returns the string representation of the provided filter flags.
func (i *filterFlags) String() string {
	var out bytes.Buffer

	for _, v := range *i {
		out.WriteString(v)
	}

	return out.String()
}

// Set sets a new filter flag in the slice.
func (i *filterFlags) Set(value string) error {
	*i = append(*i, value)

	return nil
}
