package tuna

import (
	"context"
	"net"
	"os"
	"os/signal"

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
	// graceful shutdown in interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			golog.Error("grpc: Server closed")
			g.Shutdown()
			<-g.context.Done()
		}
	}()

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
