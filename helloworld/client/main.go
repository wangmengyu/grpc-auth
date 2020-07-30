package main

import (
	"context"
	"google.golang.org/grpc"
	"grpc_auth/authentication"
	pb "grpc_auth/helloworld/helloworld"
	"log"
)

const (
	address     = "localhost:50051"
	defaultName = "wmy"
)

func main() {

	auth := authentication.Authentication{
		Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNTk2MTAwMjU0fQ.zJhx9Iiq3TwPbX1WjnwYpl8XhaXxlDcy0Q6K75dDLl8",
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithPerRPCCredentials(&auth))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx := context.Background()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: ""})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting response: %s", r.GetMessage())
}
