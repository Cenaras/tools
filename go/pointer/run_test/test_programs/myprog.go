//go:build ignore
// +build ignore

package main

// Test of context-sensitive treatment of certain function calls,
// e.g. static calls to simple accessor methods.

var a, b int

type T struct{ x *int }

func (t *T) SetX(x *int) { t.x = x }
func (t *T) GetX() *int  { return t.x }
func (t *T) Foo(x *int)  { t.SetX(x) }
func (t *T) Bar() *int   { return t.GetX() }

func context1() {
	var t1, t2 T
	t1.SetX(&a)
	t2.SetX(&b)
	x := t1.Bar()
	y := t2.Bar()
	print(x)
	print(y)
}

func main() {
	context1()
}
