package main

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"go-micro-test/proto"
)

func main() {
	//以etcd作为服务注册发现
	reg := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"etcd.sndu.cn:2379"}
	})
	//创建客户端服务
	server := micro.NewService(
		micro.Name("Hello.Client"),
		micro.Registry(reg),
	)
	//初始化
	server.Init()
	//注册要访问的方法
	userInfo := proto.NewUserServerService("Hello", server.Client())
	//调用方法
	resp, err := userInfo.UserInfo(context.Background(), &proto.GetRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Name)
	fmt.Println(resp.Age)
	fmt.Println(resp.Score)
}
