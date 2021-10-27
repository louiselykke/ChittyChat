package main

import (
	"context"
	"log"
	"net"

	"github.com/louiselykke/ChittyChat/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	proto.UnimplementedChatServer
}

func (s *Server) Brodcast(ctx context.Context, msgStream *proto.Chat_BrodcastServer) *proto.Response {
	log.Println("Brodcast call")
	msg, err := msgStream.Recv()

	for {

	}
	return status.Errorf(codes.Unimplemented, "method Publish not implemented")
}

func main() {
	// Create listener tcp on port 9080
	list, err := net.Listen("tcp", ":9080")
	if err != nil {
		log.Fatalf("Failed to listen on port 9080: %v", err)
	}
	grpcServer := grpc.NewServer()

	proto.RegisterchatServer(grpcServer, &Server{})

	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to server %v", err)
	}
}
