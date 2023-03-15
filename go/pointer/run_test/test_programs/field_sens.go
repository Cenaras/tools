//go:build ignore
// +build ignore

package main

var a, b, c int

type T struct {
	x *int
	y *int
}

func fieldSens(t1 *T, t2 *T) {
	print(t1.x)
	print(t1.y)
	print(t2.x)
	print(t2.y)
}

func main() {
	//context1()

	//var arr [2]*int
	//array5(&arr)
	var t1 *T = &T{&a, &c}
	var t2 *T = &T{&b, &a}
	fieldSens(t1, t2)
}
