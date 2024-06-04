package main

import "fmt"

type Paths []string

func (i *Paths) String() string {
	return fmt.Sprintf("%+v", *i)
}

func (i *Paths) Set(value string) error {
	*i = append(*i, value)
	return nil
}
