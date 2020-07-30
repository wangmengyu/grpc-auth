package authentication

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"time"
)

type Authentication struct {
	Token string
}
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func GetJwtKey() []byte {
	return []byte("my_secret_key")
}

// 可用的用户列表
var users = map[string]string{
	"admin": "admin",
	"test":  "test",
}

/**
  获得token
*/
func (a *Authentication) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{"token": a.Token}, nil
}

/**
  不要求底层使用安全链接
*/
func (a *Authentication) RequireTransportSecurity() bool {
	return false
}

/**
  登录的获取当前的用户名
*/
func (a *Authentication) GetCurrentUser(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing credentials")
	}
	fmt.Println("md:", md)

	var tknStr string

	if val, ok := md["token"]; ok {
		tknStr = val[0]
	}
	fmt.Printf("token=%s\n", tknStr)
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJwtKey(), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", err
		}

	}

	if !tkn.Valid {
		return "", status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return claims.Username, nil
}

func (a *Authentication) CreateToken(cred Credentials) (string, error) {
	rightPwd, ok := users[cred.Username]
	if ok == false {
		return "", fmt.Errorf("not found user")
	}

	if rightPwd != cred.Password {
		return "", fmt.Errorf("密码错误")
	}

	var creds Credentials
	// Get the JSON body and decode into credentials
	jsonStr, err := json.Marshal(cred)
	err = json.Unmarshal(jsonStr, &creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		return "", err
	}

	fmt.Println("[creds]:", creds)

	// Declare the expiration time of the token
	// here, we have kept it as 50 minutes
	expirationTime := time.Now().Add(60 * time.Minute)
	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString(GetJwtKey())
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return "", err
	}
	return tokenString, nil

}

/**
  刷新token
*/
func (a *Authentication) RefreshToken(ctx context.Context) (string, error) {

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", fmt.Errorf("missing credentials")
	}
	fmt.Println("md:", md)
	var tknStr string
	if val, ok := md["token"]; ok {
		tknStr = val[0]
	}
	fmt.Printf("token=%s\n", tknStr)

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJwtKey(), nil
	})
	if err != nil {
		return "", err
	}
	if !tkn.Valid {
		return "", fmt.Errorf("token valid error")
	}
	expirationTime := time.Now().Add(1 * time.Hour)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(GetJwtKey())
	if err != nil {
		return "", err
	}

	// Set the new token as the users `token` cookie
	return tokenString, nil

}

func (a *Authentication) Auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("missing credentials")
	}
	fmt.Println("md:", md)

	var tknStr string

	if val, ok := md["token"]; ok {
		tknStr = val[0]
	}
	fmt.Printf("token=%s\n", tknStr)
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return GetJwtKey(), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return err
		}

	}

	if tkn == nil {
		return status.Errorf(codes.Unauthenticated, "empty token")
	}

	if !tkn.Valid {
		return status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return nil
}
