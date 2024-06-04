package main

import (
	"context"
	"fmt"
	"go-micro-test/proto"
	"go-micro-test/tracer"
	"time"

	opentracingFn "github.com/go-micro/plugins/v2/wrapper/trace/opentracing"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
)

var reg = etcd.NewRegistry(func(options *registry.Options) {
	options.Addrs = []string{"10.0.33.50:2379"}
})

type UserServer struct {
	c client.Client
}

func (us *UserServer) UserInfo(ctx context.Context, res *proto.GetRequest, resp *proto.PutResponse) error {
	resp.Name = "lpj"
	resp.Age = 29
	resp.Score = 100
	time.Sleep(time.Second * 3)
	return nil
}

func (us *UserServer) UserInfoFromServer2(ctx context.Context, res *proto.GetRequest, resp *proto.PutResponse) error {
	// 从ctx中读取tracerid信息并打印
	span := opentracing.SpanFromContext(ctx)
	if span != nil {
		// 注意这里上面和下面打印的id不一致，是因为上面的转了16进制
		fmt.Printf("Reporting span %+v\n", span)
		// 获取 Trace ID
		spctx := span.Context().(jaeger.SpanContext)
		fmt.Printf("Trace ID: %d | %d | %d\n", spctx.TraceID(), spctx.ParentID(), spctx.SpanID())
	}

	time.Sleep(time.Second * 1)

	request := us.c.NewRequest("Hello.Server2", "UserServer.UserInfo", res)
	response := &proto.PutResponse{}
	err := us.c.Call(ctx, request, response)
	if err != nil {
		return err
	}
	resp.Age = response.Age
	resp.Name = response.Name
	resp.Score = response.Score
	return nil
}

func main() {
	// 链路追踪
	err := tracer.NewTracer("server", "localhost:6831")
	if err != nil {
		panic(err)
	}
	name := "Hello.Server1"

	// 做客户端
	service := micro.NewService(
		micro.Name(name),
		micro.Registry(reg),
		micro.WrapClient(opentracingFn.NewClientWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(opentracingFn.NewHandlerWrapper(opentracing.GlobalTracer())),
	)
	userService := &UserServer{}
	userService.c = service.Client()

	//服务初始化
	service.Init()
	//注册方法
	err = proto.RegisterUserServerHandler(service.Server(), userService)
	if err != nil {
		panic(err)
	}
	//启动服务
	if err = service.Run(); err != nil {
		panic(err)
	}
}
