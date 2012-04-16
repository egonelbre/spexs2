package spexs

import "fmt"

type Pattern interface {
	fmt.Stringer
}

type Reference interface {
	Next(idx int, pos byte) (Char, byte, bool)
} 

type Char rune