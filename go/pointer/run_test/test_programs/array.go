//go:build ignore
// +build ignore

package main

var a, b int

type S struct{ x *int }

func array1() {
	var x [2]*int
	x[0] = &a
	x[1] = &b
	print(x[1])
}

func array2() {
	x := new(S)
	y := &S{&a}
	c := new(int)
	*c = b
	var d *int = c
	*d = a
	print(x)
	print(y)
	print(c)
	print(d)
}

func main() {
	array1()
	array2()
}
