package tuna

import (
	"context"
	"net"

	"github.com/kataras/golog"
	"google.golang.org/grpc"
)

// GRPC is responsible for handling gRPC server.
type GRPC struct {
	server  *grpc.Server
	context context.Context
}

// Context returns the context used in gRPC.
func (g *GRPC) Context() context.Context {
	return g.context
}

// Server returns gRPC server.
func (g *GRPC) Server() *grpc.Server {
	return g.server
}

// Run starts the gRPC server on given address.
func (g *GRPC) Run(listener net.Listener) {
	// start gRPC server
	golog.Info("grpc: Server Started")
	err := g.server.Serve(listener)
	if err != nil {
		golog.Fatal(err)
	}
}

// Shutdown shutting down grpc server gracefully.
func (g *GRPC) Shutdown() {
	g.server.GracefulStop()
}

// New instantiate the GRPC handler.
func New(ctx context.Context, opt ...grpc.ServerOption) *GRPC {
	opt = append(opt, grpc.UnaryInterceptor(serverLogInterceptor))

	return &GRPC{
		server:  grpc.NewServer(opt...),
		context: ctx,
	}
}
