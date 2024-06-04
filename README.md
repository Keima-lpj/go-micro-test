# go-micro-test

# 先启动jaeger

docker run -d -p 6831:6831/udp -p 16686:16686 jaegertracing/all-in-one:latest

# 然后依次启动server2、server、运行client，观察链路

# 注意，在jaeger web服务中需要使用traceid的16进制格式搜索