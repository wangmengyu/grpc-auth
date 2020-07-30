package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"grpc_auth/authentication"
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

func backup() {
	auth := authentication.Authentication{
		User:     "wmy",
		Password: "555",
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "wmy"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting response: %s", r.GetMessage())
}
