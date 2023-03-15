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
	x := t1.GetX()
	y := t2.GetX()
	print(x)
	print(y)
}

func main() {
	//context1()

	//var arr [2]*int
	//array5(&arr)
}

func array5(arr *[2]*int) {
	var x int
	arr[0] = &x
	arr[1] = &b

	var n int
	print(arr[n]) // @pointsto command-line-arguments.a | command-line-arguments.b
}
