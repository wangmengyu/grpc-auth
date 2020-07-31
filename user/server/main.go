package main

import (
	"context"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"grpc_auth/authentication"
	"grpc_auth/interceptor"
	"grpc_auth/user/user"
	"log"
	"net"
)

const (
	port = ":50052"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	auth       *authentication.Authentication
	authFilter *interceptor.AuthFilter
}

//获取当前用户的信息
func (s *server) Info(ctx context.Context, req *user.InfoReq) (*user.InfoRep, error) {
	username, err := s.auth.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	pwd, ok := authentication.GetPwdByUser(username)
	if !ok {
		return nil, fmt.Errorf("not found user")
	}

	return &user.InfoRep{
		Username: username,
		Password: pwd,
	}, nil

}

/**
用户列表
*/
func (s *server) List(req *user.ListReq, stream user.User_ListServer) error {
	userList := authentication.GetUserList()
	for username, pwd := range userList {
		item := new(user.InfoRep)
		item.Username = username
		item.Password = pwd
		if err := stream.Send(item); err != nil {
			return err
		}
	}
	return nil
}

func main() {

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	af := new(interceptor.AuthFilter)
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			af.AuthUnary,
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			af.AuthStream,
		)),
	)
	user.RegisterUserServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
