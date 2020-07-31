cd /Users/mengyuwang/work/code/grpc-auth/helloworld
protoc -I ./ ./helloworld/helloworld.proto --go_out=plugins=grpc:helloworld


cd /Users/mengyuwang/work/code/grpc-auth/user
protoc -I ./ ./user/user.proto --go_out=plugins=grpc:user