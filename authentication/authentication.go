package authentication

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Authentication struct {
	User     string
	Password string
}

/**
  获得用户名和密码
*/
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"user": a.User, "password": a.Password}, nil
}

/**
  不要求底层使用安全链接
*/
func (a *Authentication) RequireTransportSecurity() bool {
	return false
}

func (a *Authentication) Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}
	fmt.Println("md:", md)

	var appid string
	var appkey string

	if val, ok := md["user"]; ok {
		appid = val[0]
	}
	if val, ok := md["password"]; ok {
		appkey = val[0]
	}

	fmt.Printf("appid=%s, appkey=%s\n", appid, appkey)

	if appid != "admin" || appkey != "admin" {
		// 此处要验证key的正确性
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return nil
}
