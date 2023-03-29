package main

import _ "net"

func main() {
	ch := make(chan int)
	<-ch //@ analysis(false)

}
