package main

import (
	"context"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"go-micro-test/proto"
)

type UserServer struct{}

func (us *UserServer) UserInfo(ctx context.Context, res *proto.GetRequest, resp *proto.PutResponse) error {
	resp.Name = "lpj"
	resp.Age = 29
	resp.Score = 100
	return nil
}

func main() {
	//以etcd作为服务注册发现
	reg := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"etcd.sndu.cn:2379"}
	})
	//创建一个新的服务
	server := micro.NewService(
		micro.Name("Hello"),
		micro.Registry(reg),
	)
	//服务初始化
	server.Init()
	//注册方法
	err := proto.RegisterUserServerHandler(server.Server(), new(UserServer))
	if err != nil {
		panic(err)
	}
	//启动服务
	if err = server.Run(); err != nil {
		panic(err)
	}
}
