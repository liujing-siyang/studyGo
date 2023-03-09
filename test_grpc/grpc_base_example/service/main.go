package main

import (
	"context"
	"fmt"
	pb "baseexample/pb"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// hello server

var blacklist = map[string]int{
	"liming": 1,
	"world":  1,
}

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	_, ok := blacklist[in.Name]
	if ok {
		st := status.New(codes.ResourceExhausted, "Request Obeject in blacklist.")
		ds, err := st.WithDetails(
			&errdetails.QuotaFailure{
				Violations: []*errdetails.QuotaFailure_Violation{{
					Subject:     fmt.Sprintf("name:%s", in.Name),
					Description: "该用户属于黑名单用户",
				}},
			},
		)
		if err != nil {
			return nil, st.Err()
		}
		return nil, ds.Err()
	}
	return &pb.HelloResponse{Reply: "Hello " + in.Name}, nil
}

// LotsOfReplies 返回使用多种语言打招呼
func (s *server) LotsOfReplies(in *pb.HelloRequest, stream pb.Greeter_LotsOfRepliesServer) error {
	words := []string{
		"你好",
		"hello",
		"こんにちは",
		"안녕하세요",
	}

	for _, word := range words {
		data := &pb.HelloResponse{
			Reply: word + in.GetName(),
		}
		// 使用Send方法返回多个数据
		if err := stream.Send(data); err != nil {
			return err
		}
	}
	return nil
}

// LotsOfGreetings 接收流式数据
func (s *server) LotsOfGreetings(stream pb.Greeter_LotsOfGreetingsServer) error {
	reply := "你好："
	for {
		// 接收客户端发来的流式数据
		res, err := stream.Recv()
		if err == io.EOF {
			// 最终统一回复
			return stream.SendAndClose(&pb.HelloResponse{
				Reply: reply,
			})
		}
		if err != nil {
			return err
		}
		reply += res.GetName()
	}
}

// BidiHello 双向流式打招呼
func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	for {
		// 接收流式请求
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		reply := magic(in.GetName()) // 对收到的数据做些处理

		// 返回流式响应
		if err := stream.Send(&pb.HelloResponse{Reply: reply}); err != nil {
			return err
		}
	}
}

// magic 一段价值连城的“人工智能”代码
func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
}

// UnarySayHello 普通RPC调用服务端metadata操作
func (s *server) UnarySayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	// 通过defer中设置trailer.
	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		grpc.SetTrailer(ctx, trailer)
	}()

	// 从客户端请求上下文中读取metadata.
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.DataLoss, "UnarySayHello: failed to get metadata")
	}
	if t, ok := md["token"]; ok {
		fmt.Printf("token from metadata:\n")
		if len(t) < 1 || t[0] != "app-test-q1mi" {
			return nil, status.Error(codes.Unauthenticated, "认证失败")
		}
	}

	// 创建和发送header.
	header := metadata.New(map[string]string{"location": "BeiJing"})
	grpc.SendHeader(ctx, header)

	fmt.Printf("request received: %v, say hello...\n", in)

	return &pb.HelloResponse{Reply: in.Name}, nil
}

// BidirectionalStreamingSayHello 流式RPC调用客户端metadata操作
func (s *server) BidirectionalStreamingSayHello(stream pb.Greeter_BidirectionalStreamingSayHelloServer) error {
	// 在defer中创建trailer记录函数的返回时间.
	defer func() {
		trailer := metadata.Pairs("timestamp", strconv.Itoa(int(time.Now().Unix())))
		stream.SetTrailer(trailer)
	}()

	// 从client读取metadata.
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.DataLoss, "BidirectionalStreamingSayHello: failed to get metadata")
	}

	if t, ok := md["token"]; ok {
		fmt.Printf("token from metadata:\n")
		for i, e := range t {
			fmt.Printf(" %d. %s\n", i, e)
		}
	}

	// 创建和发送header.
	header := metadata.New(map[string]string{"location": "X2Q"})
	stream.SendHeader(header)

	// 读取请求数据发送响应数据.
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Printf("request received %v, sending reply\n", in)
		if err := stream.Send(&pb.HelloResponse{Reply: in.Name}); err != nil {
			return err
		}
	}
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
