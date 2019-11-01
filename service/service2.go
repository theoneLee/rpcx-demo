package main

import (
	"flag"
	"fmt"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"log"
	"rpcx-demo/service/product"
	"rpcx-demo/service/user"
	"time"
)

var (
	addr2 = flag.String("addr2", "localhost:8973", "server address 2")

	//etcd
	basePath = flag.String("base", "/rpcx_test", "prefix path")
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
)

func main() {
	flag.Parse()

	s := server.NewServer()

	addRegistryPlugin2(s) //etcd

	s.RegisterName("ProductImage", product.New("./product/static"), "")
	s.RegisterName("auth", user.New(), "") //user.New()返回一个服务对象，该服务对象的所有方法都会是允许rpc调用的，只要符合方法签名
	err := s.Serve("tcp", *addr2)          //todo 运行第二个不同端口的server
	if err != nil {
		panic(err)
	}
}

func addRegistryPlugin2(s *server.Server) {
	fmt.Println("*addr:", *addr2)
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + *addr2,
		EtcdServers:    []string{*etcdAddr},
		BasePath:       *basePath,
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		log.Fatal(err)
	}
	s.Plugins.Add(r)
}
