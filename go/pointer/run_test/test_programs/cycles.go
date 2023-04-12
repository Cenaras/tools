//go:build ignore
// +build ignore

package main

type S struct{ x int }

var c S

func cycles1() {
	a := &c
	d := c
	b := *a
	*a = b

	print(a.x)
	print(b)
	print(c.x)
	print(d)
}

func cycles2() {
	x := S{2}
	y := S{3}
	x, y = y, x
}

func main() {
	//cycles1()
	cycles2()
}
