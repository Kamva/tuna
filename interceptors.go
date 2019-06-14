package tuna

import (
	"context"

	"github.com/kataras/golog"
	"google.golang.org/grpc"
)

func serverLogInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	h, err := handler(ctx, req)

	if err != nil {
		golog.Errorf("grpc: %s [err: %s]", info.FullMethod, err)
	} else {
		golog.Infof("grpc: %s [OK]", info.FullMethod)
	}

	return h, err
}
