package grpc

import (
	"advertisingService/internal/adApp"
	"context"
	"google.golang.org/grpc"
	"time"

	"advertisingService/internal/userApp"
	"log"
	"net"
)

type Server struct {
	srv *grpc.Server
	lis net.Listener
}

func Logger(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	log.Printf("Request - Method: %s\tDuration: %s\tError: %v\n", info.FullMethod, time.Since(start), err)

	return h, err
}

func Panic(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("Request - Method: %s\tError: %v\n\n", info.FullMethod, err)
		}
	}()

	return handler(ctx, req)
}

func NewGRPCServer(lis net.Listener, a adApp.App, ua userApp.App) Server {
	srv := grpc.NewServer(grpc.ChainUnaryInterceptor(Logger), grpc.ChainUnaryInterceptor(Panic))
	RegisterAdServiceServer(srv, NewServiceServer(a, ua))
	s := Server{srv: srv, lis: lis}
	return s
}

func (s *Server) SetListener(lis net.Listener) {
	s.lis = lis
}

func (s *Server) Listen() error {
	return s.srv.Serve(s.lis)
}

func (s *Server) Stop() {
	s.srv.Stop()
}

func (s *Server) GetServer() *grpc.Server {
	return s.srv
}
