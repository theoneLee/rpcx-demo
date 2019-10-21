package main

import (
	"flag"
	"rpcx-demo/service/product"
	"rpcx-demo/service/user"

	"github.com/smallnest/rpcx/server"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")
)

func main() {
	flag.Parse()

	s := server.NewServer()
	s.RegisterName("ProductImage", product.New("./product/static"), "")
	s.RegisterName("auth", user.New(), "")
	err := s.Serve("tcp", *addr)
	if err != nil {
		panic(err)
	}
}
