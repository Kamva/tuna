package tuna

import (
	"context"
	"net"
	"os"
	"os/signal"

	"github.com/Kamva/shark"
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
func (g *GRPC) Run(address string, ch chan byte) {
	listen, err := net.Listen("tcp", address)
	shark.PanicIfError(err)

	// graceful shutdown in interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			golog.Error("grpc: Server closed")

			g.server.GracefulStop()

			<-g.context.Done()
		}
	}()

	// start gRPC server
	golog.Infof("grpc: Server Started on %s", address)
	shark.PanicIfError(g.server.Serve(listen))

	ch <- 1
}

// New instantiate the GRPC handler.
func New(opt ...grpc.ServerOption) *GRPC {
	opt = append(opt, grpc.UnaryInterceptor(serverLogInterceptor))

	return &GRPC{
		server:  grpc.NewServer(opt...),
		context: context.Background(),
	}
}
