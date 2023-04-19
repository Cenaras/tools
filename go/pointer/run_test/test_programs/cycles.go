//go:build ignore
// +build ignore

package main

type I interface {
	foo() *int
}

type S struct{ x *int }

func (s *S) foo() *int {
	c = s.x
	d = &a
	e = &d
	*e = c
	c = *e
	return *e
}

var a, b int
var c, d *int
var e **int

func cycles1() {
	c = &a
	d = &a
	e = &d
	*e = c
	c = *e
	print(c)
}

func cycles2() {
	var s I = &S{&b}
	t := s.foo()
	print(t)
}

func main() {
	//cycles1()
	cycles2()
}
