package services

import (
	"context"

	gw "dig-inv/gen/go"
	"google.golang.org/grpc/grpclog"
)

// Implements of EchoServiceServer

type echoServer struct {
	gw.UnimplementedYourServiceServer
}

func NewEchoServer() gw.YourServiceServer {
	return new(echoServer)
}

func (s *echoServer) Echo(ctx context.Context, msg *gw.StringMessage) (*gw.StringMessage, error) {
	grpclog.Info(msg)
	return msg, nil
}
