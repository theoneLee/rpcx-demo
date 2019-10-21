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
	s.RegisterName("auth", user.New(), "") //user.New()返回一个服务对象，该服务对象的所有方法都会是允许rpc调用的，只要符合方法签名
	err := s.Serve("tcp", *addr)
	if err != nil {
		panic(err)
	}
}
