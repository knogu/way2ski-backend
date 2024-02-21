package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "way2ski-backend/proto"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedWayServiceServer
}

func (s *server) GetLines(ctx context.Context, in *pb.Params) (*pb.Lines, error) {
	return &pb.Lines{
		Lines: nil,
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterWayServiceServer(s, &server{})
	reflection.Register(s)
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
