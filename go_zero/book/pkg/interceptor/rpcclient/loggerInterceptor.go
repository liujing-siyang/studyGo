package rpcclient

import (
	"context"
	"fmt"
	"path"

	"github.com/zeromicro/go-zero/core/breaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// rpc client logger interceptor
func ClientLoggerInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 通过metadata在服务之间传递数据
	md := metadata.New(map[string]string{"username": "zhangsan"})
	ctx = metadata.NewOutgoingContext(ctx, md)

	fmt.Println("拦截前")
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		return err
	}
	fmt.Println("拦截后")
	return nil
}

func BreakerInterceptor(ctx context.Context, method string, req, reply interface{},
	cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 基于请求方法进行熔断
	breakerName := path.Join(cc.Target(), method)
	return breaker.DoWithAcceptable(breakerName, func() error {
		// 真正发起调用
		return invoker(ctx, method, req, reply, cc, opts...)
		// DefaultAcceptable判断哪种错误需要加入熔断错误计数
	}, DefaultAcceptable)
}

func DefaultAcceptable(err error) bool {
	return err == nil
}
