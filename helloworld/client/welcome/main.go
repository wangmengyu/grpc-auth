package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "grpc_auth/helloworld/helloworld"
	"log"
)

var ctx context.Context

func init() {
	ctx = context.Background()
}

const (
	address = "localhost:50051"
)

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
	req := new(pb.WelReq)
	req.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNTk2MDk1NzIzfQ.rR3Becovrp_J3tyKwufqbDDApMKTY7-8c18UMZMB0xg"

	res, err := client.Welcome(ctx, req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println(res)
}
