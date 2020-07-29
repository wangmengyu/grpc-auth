package authentication

import "context"

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
