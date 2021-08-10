package main

import (
	"context"
	"fmt"
	"github.com/xhyonline/grpc-graceful-stop/gen"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := gen.NewGracefulClient(conn)

	// 100 个并发请求
	pool:=make(chan struct{},100)
	var count int
	for  {
		pool<- struct{}{}
		count++
		go func(count int) {
			result, err := client.GracefulMessage(context.Background(),&gen.Request{
				Name: "小明",
				Age:  24,
			})
			if err != nil {
				log.Fatalf("%s", err)
			}
			fmt.Printf("第%d个请求:%+v\n", count,result.GetSelfDescription())
			<-pool
		}(count)
	}

}

