package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"grpc_auth/authentication"
	pb "grpc_auth/helloworld/helloworld"
	"log"
	"net"
	"time"
)

var jwtKey = []byte("my_secret_key")

// 可用的用户列表
var users = map[string]string{
	"admin": "admin",
	"test":  "test",
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

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
	return &pb.HelloReply{Message: "Hello " + in.GetName() + "!"}, nil
}
func (s *server) Login(ctx context.Context, in *pb.LoginReq) (*pb.LoginRep, error) {

	rightPwd, ok := users[in.Username]
	if ok == false {
		return nil, fmt.Errorf("not found user")
	}

	if rightPwd != in.Password {
		return nil, fmt.Errorf("密码错误")
	}

	var creds Credentials
	// Get the JSON body and decode into credentials
	jsonStr, err := json.Marshal(in)
	err = json.Unmarshal(jsonStr, &creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		return nil, err
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
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		// If there is an error in creating the JWT return an internal server error
		return nil, err
	}

	// Finally, we set the client cookie for "token" as the JWT we just generated
	// we also set an expiry time which is the same as the token itself
	return &pb.LoginRep{Token: tokenString}, nil
}

func (s *server) Welcome(ctx context.Context, in *pb.WelReq) (*pb.WelRep, error) {

	tknStr := in.Token

	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}

	}
	if !tkn.Valid {
		return nil, fmt.Errorf("valid token")
	}

	fmt.Printf("Welcome %s!\n", claims.Username)
	return new(pb.WelRep), nil

}

func (s *server) Refresh(ctx context.Context, in *pb.RefReq) (*pb.RefRep, error) {

	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route

	tknStr := in.GetToken()
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !tkn.Valid {
		return nil, fmt.Errorf("token valid error")
	}
	// (END) The code up-till this point is the same as the first part of the `Welcome` route

	// Now, create a new token for the current use, with a renewed expiration time
	expirationTime := time.Now().Add(1 * time.Hour)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, err
	}

	// Set the new token as the users `token` cookie
	return &pb.RefRep{
		Token: tokenString,
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
