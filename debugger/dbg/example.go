package main

import (
	"fmt"

	"github.com/egonelbre/spexs2/debugger"
)

type BinOp func(int, int) int

func Adder(a int, b int) int {
	return a + b
}

func AttachDebugger(d *debugger.Debugger, f BinOp) BinOp {
	return func(a int, b int) int {
		result := f(a, b)
		d.Break(func() {
			fmt.Fprintf(d.Logout, "f(%v,%v) = %v\n", a, b, result)
		})
		return result
	}
}

func main() {
	d := debugger.New()
	debugger.RunShell(d)

	adder := AttachDebugger(d, Adder)
	done := make(chan int, 100)

	for i := 0; i < 10; i++ {
		go func(a int) {
			fmt.Printf(".")
			adder(a, a*3)
			fmt.Printf("-")
			done <- 1
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	print("\n\n")
}
