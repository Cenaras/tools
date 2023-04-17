//go:build ignore
// +build ignore

package main

var a, b int

type I interface {
	foo(x *int) *S
	bar(x *int) *S
}

type T struct{}

type S struct {
	x *int
}

func (t *T) bar(x *int) *S {
	return &S{x}
}

var s1 *S

func heap1() {
	var t1 = &T{}
	var t2 = &T{}
	s1 = t1.bar(&a)
	s1 = t2.bar(&b)
	print(s1)
	//print(s2)
}

func (t *T) foo(x *int) *S {
	var k I = &T{}
	return k.bar(x)
}

func heap2() {
	var t1 I = &T{}
	var t2 I = &T{}
	s1 := t1.foo(&a)
	s2 := t2.foo(&b)
	print(s1.x)
	print(s2.x)
}

func main() {
	heap1()
	//heap2()
}
