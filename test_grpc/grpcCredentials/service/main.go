package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	pb "gprccredentials/pb"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// hello server

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

func main() {
	// 监听本地的8972端口
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	//添加证书,单项认证
	// creds,err :=  credentials.NewServerTLSFromFile("cert/server.pem","cert/server.key")
	// if err != nil{
	// 	log.Fatal("证书生成错误",err)
	// }

	//双向认证
	// 从证书相关文件中读取和解析信息，得到证书公钥，密钥对
	cert, err := tls.LoadX509KeyPair("../cert/server.pem", "../cert/server.key")
	if err != nil {
		log.Fatal("证书生成错误", err)
	}
	certPool := x509.NewCertPool()
	ca,err := ioutil.ReadFile("../cert/ca.crt")
	if err != nil {
		log.Fatal("证书生成错误", err)
	}
	//尝试解析所传入的PEM编码证书
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		//设置证书链
		Certificates: []tls.Certificate{cert},
		//要求必须校验客户端的证书
		ClientAuth: tls.RequireAndVerifyClientCert,
		//设置根证书集合，校验方式使用ClientAuth中设定的模式
		ClientCAs: certPool,
	})
	s := grpc.NewServer(grpc.Creds(creds)) // 创建gRPC服务器
	pb.RegisterGreeterServer(s, &server{}) // 在gRPC服务端注册服务
	// 启动服务
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
