package interceptor

import (
	"context"
	"google.golang.org/grpc"
	"grpc_auth/authentication"
	"log"
)

type AuthFilter struct {
}

/**
实现普通方法的截取器,
需要为grpc.UnaryInterceptor的参数实现一个函数：
函数的ctx和req参数就是每个普通的RPC方法的前两个参数。
 第三个info参数表示当前是对应的那个gRPC方法，第四个handler参数对应当前的gRPC方法函数。
函数中首先是日志输出info参数，然后调用handler对应的gRPC方法函数。
要使用filter截取器函数，只需要在启动gRPC服务时作为参数输入即可：
*/
func (af *AuthFilter) AuthUnary(ctx context.Context,
	req interface{}, info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	log.Println("auth-filter:", info) // 先是日志输出info参数
	//然后做验证
	auth := new(authentication.Authentication)
	if err := auth.Auth(ctx); err != nil {
		return nil, err
	}
	return handler(ctx, req) // 然后调用handler对应的gRPC方法函数v
}

func (af *AuthFilter) AuthStream(srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	auth := new(authentication.Authentication)
	if err := auth.Auth(stream.Context()); err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println("--> stream interceptor: ", info.FullMethod)
	return handler(srv, stream)
}
