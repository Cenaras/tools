//go:build ignore
// +build ignore

package main

var a, b int

type I interface {
	foo() *int
}

type S struct{}

func (s *S) foo() *int {
	return &a
}

func staticInvoke1() {
	var s *S
	var x *int
	x = s.foo()
	print(x)
}

func staticInvoke2() {
	var s *S
	var i I = s
	print(i.foo())
}

func staticInvoke3() {
	var i I
	print(i.foo())
}

func main() {
	staticInvoke1()
	staticInvoke2()
	staticInvoke3()
}
