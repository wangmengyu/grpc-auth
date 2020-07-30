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
	req := new(pb.RefReq)
	req.Token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNTk2MDk1Njg4fQ.d9s0vUxqGVrDtfkURXX5oEMgftFqZmKmloRuVgqRb58"

	res, err := client.Refresh(ctx, req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println(res.Token)
}