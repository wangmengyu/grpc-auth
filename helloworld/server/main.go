package main

import (
	"context"
	"google.golang.org/grpc"
	"grpc_auth/authentication"
	pb "grpc_auth/helloworld/helloworld"
	"log"
	"net"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	auth *authentication.Authentication
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	if err := s.auth.Auth(ctx); err != nil {
		return nil, err
	}
	user, err := s.auth.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	return &pb.HelloReply{Message: "Hello " + user + "!"}, nil
}

/**
登录
*/
func (s *server) Login(ctx context.Context, in *pb.LoginReq) (*pb.LoginRep, error) {
	var cred authentication.Credentials
	cred.Username = in.Username
	cred.Password = in.Password
	tokenString, err := s.auth.CreateToken(cred)
	if err != nil {
		return nil, err
	}
	return &pb.LoginRep{Token: tokenString}, nil
}

func (s *server) Refresh(ctx context.Context, in *pb.RefReq) (*pb.RefRep, error) {

	// 检查当前登录状态
	if err := s.auth.Auth(ctx); err != nil {
		return nil, err
	}
	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route

	t, err := s.auth.RefreshToken(ctx)

	if err != nil {
		return nil, err
	}

	return &pb.RefRep{
		Token: t,
	}, nil

}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
