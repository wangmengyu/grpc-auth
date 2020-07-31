package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc_auth/authentication"
	"grpc_auth/user/user"
	"io"
	"log"
)

var ctx context.Context

func init() {
	ctx = context.Background()
}

const (
	address = "localhost:50052"
)

func main() {
	// 连接 rpc server
	auth := authentication.Authentication{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNTk2MTY2NjA5fQ.KRbp1BQssJAVJIwnZqqP2RxVLA0NTH3Ih5EXcgGGrPg",
	}
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer conn.Close()
	client := user.NewUserClient(conn)

	// 保存到本地数据库
	req := new(user.ListReq)

	stream, err := client.List(ctx, req)
	if err != nil {
		log.Println(err.Error())
		return
	}

	items := make([]*user.InfoRep, 0)

	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err.Error())
			return
		}
		if item != nil {
			items = append(items, item)
		}
	}

	if err != nil {
		log.Println(err.Error())
		return
	}

	fmt.Println(items)
}
