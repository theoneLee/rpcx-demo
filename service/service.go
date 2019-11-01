package main

import (
	"flag"
	"fmt"
	"github.com/smallnest/rpcx/serverplugin"
	"log"
	"rpcx-demo/service/product"
	"rpcx-demo/service/user"
	"time"

	"github.com/smallnest/rpcx/server"
)

var (
	addr = flag.String("addr", "localhost:8972", "server address")

	//etcd
	basePath = flag.String("base", "/rpcx_test", "prefix path")
	etcdAddr = flag.String("etcdAddr", "localhost:2379", "etcd address")
)

func main() {
	flag.Parse()

	s := server.NewServer()

	addRegistryPlugin(s) //etcd

	s.RegisterName("ProductImage", product.New("./product/static"), "")
	s.RegisterName("auth", user.New(), "") //user.New()返回一个服务对象，该服务对象的所有方法都会是允许rpc调用的，只要符合方法签名
	err := s.Serve("tcp", *addr)
	if err != nil {
		panic(err)
	}
}

//服务发现：可以感知集群上哪些提供了同一种服务（节点上线，节点下线会被etcd感知，然后客户端利用负载均衡规则时会选择好节点），
// 然后利用负载均衡规则选择一个节点给客户端调用。。比如可以运行几个不同端口的server作为集群，然后给客户端调用。
// todo 有能力可以结合原生的etcd实现服务发现的内容看一下，他是怎么利用watch实现服务发现的
//https://github.com/etcd-io/etcd/tree/master/clientv3   etcd客户端文档
//https://github.com/daizuozhuo/etcd-service-discovery
//https://github.com/bailu1901/pkg/blob/master/etcd_watcher/watcher_test.go
func addRegistryPlugin(s *server.Server) {
	fmt.Println("*addr:", *addr)
	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: "tcp@" + *addr,
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
