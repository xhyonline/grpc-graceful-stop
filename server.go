package main

import (
	"context"
	"fmt"
	"github.com/xhyonline/grpc-graceful-stop/gen"
	"github.com/xhyonline/xutil/sig"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

var g = grpc.NewServer()

type Message struct {
}

func (s *Message) GracefulMessage(_ context.Context, request *gen.Request) (*gen.Response, error) {
	resp := &gen.Response{
		SelfDescription: fmt.Sprintf("我的名字是%s,今年%d岁", request.GetName(), request.GetAge()),
	}
	time.Sleep(time.Second * 5)
	return resp, nil
}

// GracefulClose 优雅停止
func (s *Message) GracefulClose() {
	fmt.Println("开始进行优雅停止")
	g.GracefulStop()
}

func main() {
	m := new(Message)
	gen.RegisterGracefulServer(g, m)
	l, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatal(err)
	}
	ctx := sig.Get().RegisterClose(m)
	// 服务启动
	go g.Serve(l)
	select {
		case <-ctx.Done():
			fmt.Println("grpc 服务优雅退出完毕")
	}
}
