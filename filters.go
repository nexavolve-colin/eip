package main

import (
	"bytes"
)

type filterFlags []string

func (i *filterFlags) String() string {
	var out bytes.Buffer

	for _, v := range *i {
		out.WriteString(v)
	}

	return out.String()
}

func (i *filterFlags) Set(value string) error {
	*i = append(*i, value)

	return nil
}
