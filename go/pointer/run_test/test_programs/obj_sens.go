//go:build ignore
// +build ignore

package main

// Test of context-sensitive treatment of certain function calls,
// e.g. static calls to simple accessor methods.

var a, b int

var unknown bool

func id(x *int) *int {
	return x
}

type I interface {
	foo(x *int) *int
	baz(x *int) *int
	bar()
}

type T struct{ x *int }

func (t *T) foo(x *int) *int {
	return x
}

type S struct {
	x *int
	y *int
}

func (s *S) foo(x *int) *int {
	return x
}

func (s *S) baz(x *int) *int {
	var i I = &S{}
	return i.foo(x)
}

func (s *S) bar() {
	s.y = s.x
}

func context1() {
	var s I = &S{}
	print(s.foo(&a))
	print(s.foo(&b))
}

func context2() {
	var s *S
	s1 := &S{x: &a}
	s2 := &S{x: &b}
	if unknown {
		s = s1
	} else {
		s = s2
	}
	baz(s)
	//s1.bar()
	print(s1.y)
	print(s2.y)
	//print(s.foo(&a))
	//print(s.foo(&b))
	//print(s1.foo(&b))
}

func baz(i I) {
	i.bar()
}

func context3() {
	var s1 I = &S{}
	var s2 I = &S{}
	print(s1.baz(&a))
	print(s2.baz(&b))
}

func main() {
	context1()
	context2()
	context3()
}
