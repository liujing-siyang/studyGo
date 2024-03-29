package main

import (
	"context"
	"fmt"
	pb "grpcimport/pb"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/anypb"
)

// hello server

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	user := pb.User{Username: "lihua",Age: 18}
	content := pb.Content{Msg: "msg ..."}
	//转换成any类型
	any,_ := anypb.New(&content)
	return &pb.HelloResponse{Reply: "Hello " + in.Name,User: &user,Data: any}, nil
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	s := grpc.NewServer()                  // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server{}) // 在gRPC服务端注册服务
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
