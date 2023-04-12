//go:build ignore
// +build ignore

package main

type S struct{ x *int }

var a, b int

var c, d *int

func cycles1() {
	c = &a
	d = &b
	var e = &d
	*e = c
	c = *e

	print(c)
}

func cycles2() {

}

func main() {
	cycles1()
	//cycles2()
}
