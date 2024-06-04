package main

import (
	"context"
	"fmt"
	"go-micro-test/proto"
	"go-micro-test/tracer"

	opentracingFn "github.com/go-micro/plugins/v2/wrapper/trace/opentracing"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/opentracing/opentracing-go"
)

func main() {
	//以etcd作为服务注册发现
	reg := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"10.0.33.50:2379"}
	})
	// 链路追踪
	err := tracer.NewTracer("client", "localhost:6831")
	if err != nil {
		panic(err)
	}
	//创建客户端服务
	server := micro.NewService(
		micro.Name("client"),
		micro.Registry(reg),
		micro.WrapClient(opentracingFn.NewClientWrapper(opentracing.GlobalTracer())),
	)
	//初始化
	server.Init()
	//注册要访问的方法
	userInfo := proto.NewUserServerService("Hello.Server1", server.Client())
	//调用方法
	resp, err := userInfo.UserInfoFromServer2(context.Background(), &proto.GetRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.Name)
	fmt.Println(resp.Age)
	fmt.Println(resp.Score)
	select {}
}
