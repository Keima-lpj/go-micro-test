package main

import (
	"context"
	"fmt"
	"go-micro-test/proto"
	"go-micro-test/tracer"
	"time"

	opentracingFn "github.com/go-micro/plugins/v2/wrapper/trace/opentracing"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

type UserServer struct{}

func (us *UserServer) UserInfo(ctx context.Context, res *proto.GetRequest, resp *proto.PutResponse) error {
	// 从ctx中读取tracerid信息并打印
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		// 获取 Trace ID
		traceID := span.Context().(jaeger.SpanContext).TraceID()
		fmt.Printf("Trace ID: %d\n", traceID)
	}

	resp.Name = "lpj1"
	resp.Age = 32
	resp.Score = 100
	time.Sleep(time.Second * 2)
	return nil
}

func (us *UserServer) UserInfoFromServer2(ctx context.Context, res *proto.GetRequest, resp *proto.PutResponse) error {
	panic(111)
}

func main() {
	//以etcd作为服务注册发现
	reg := etcd.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{"10.0.33.50:2379"}
	})
	// 链路追踪
	err := tracer.NewTracer("server2", "localhost:6831")
	if err != nil {
		panic(err)
	}
	//创建一个新的服务
	server := micro.NewService(
		micro.Name("Hello.Server2"),
		micro.Registry(reg),
		micro.WrapHandler(opentracingFn.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	//服务初始化
	server.Init()
	//注册方法
	err = proto.RegisterUserServerHandler(server.Server(), new(UserServer))
	if err != nil {
		panic(err)
	}
	//启动服务
	if err = server.Run(); err != nil {
		panic(err)
	}
}
