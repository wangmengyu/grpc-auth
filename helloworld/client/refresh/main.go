package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc_auth/authentication"
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
	auth := authentication.Authentication{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNTk2MTAyMDkwfQ.SJzFlFg7K_Q-JelHUi342LUFVJlRsV4enSPTypoOx48",
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer conn.Close()
	client := pb.NewGreeterClient(conn)

	// 保存到本地数据库
	req := new(pb.RefReq)

	res, err := client.Refresh(ctx, req)

	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println(res.Token)
}
