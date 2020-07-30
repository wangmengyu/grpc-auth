package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc_auth/helloworld/helloworld"
	"log"
)

const (
	address     = "localhost:50051"
	defaultName = "wmy"
)

var ctx context.Context

func init() {
	ctx = context.Background()
}

func main() {
	// 连接 rpc server
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// 保存到本地数据库
	req := new(pb.LoginReq)
	req.Username = "admin"
	req.Password = "admin"
	res, err := client.Login(ctx, req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println(res)

}
