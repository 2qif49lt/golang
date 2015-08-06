package main

import (
	"fmt"
)

type cb interface {
	add(a, b int) int
}

func caller(c cb) {
	fmt.Println(c.add(1, 1))
}

type cbwarpper func(a, b int) int

func (c cbwarpper) add(a, b int) int {
	return c(a, b)
}

type decorator func(cb) cb

func decorate(last cb, decors ...decorator) cb {
	for _, decor := range decors {
		last = decor(last)
	}
	return last
}

func logpara() decorator {
	return func(c cb) cb {
		return cbwarpper(func(a, b int) int {
			fmt.Println("before")
			return c.add(a, b)
		})
	}
}

func main() {
	caller(decorate(cbwarpper(func(a, b int) int {
		return a + b + 1
	}), logpara()))
}
