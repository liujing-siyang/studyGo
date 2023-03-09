package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"time"

	pb "gprccredentials/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// hello_client

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "127.0.0.1:8972", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	//添加公钥，单项认证
	// creds, err := credentials.NewClientTLSFromFile("cert/server.pem", "*.jyw.com")
	// if err != nil {
	// 	log.Fatal("证书生成错误", err)
	// }
	cert, err := tls.LoadX509KeyPair("../cert/client.pem", "../cert/client.key")
	if err != nil {
		log.Fatal("证书生成错误", err)
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("../cert/ca.crt")
	if err != nil {
		log.Fatal("证书生成错误", err)
	}
	//尝试解析所传入的PEM编码证书
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		//设置证书链
		Certificates: []tls.Certificate{cert},
		//通过的域名
		ServerName: "*.jyw.com",
		//设置根证书集合，校验方式使用ClientAuth中设定的模式
		RootCAs: certPool,
	})
	// 连接到server端，使用证书
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetReply())
}
