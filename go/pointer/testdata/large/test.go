package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
)
func main() {
	fmt.Println("Hello")
	fmt.Println(strings.ToUpper("hi"))
	conn, _ := net.Dial("tcp", "golang.org:80")
	print(conn)
	print(os.Getenv("ETCD_CLIENT_DEBUG"))
	print(time.Duration(0))
	print(context.Canceled)
	print(uuid.New())
}