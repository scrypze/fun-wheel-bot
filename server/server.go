package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"math/rand"

	"google.golang.org/grpc"
	pb "fun-wheel-bot/grpc" 
)

const port = ":50051"

type server struct {
    pb.UnimplementedFunWheelServiceServer
    wheels map[int64][]string
}

func (s *server) CreateWheel(ctx context.Context, req *pb.CreateWheelRequest) (*pb.CreateWheelResponse, error) {
    s.wheels[req.GetChatId()] = []string{}
    return &pb.CreateWheelResponse{Message: "Wheel created!"}, nil
}

func (s *server) AddOption(ctx context.Context, req *pb.AddOptionRequest) (*pb.AddOptionsResponse, error) {
    s.wheels[req.GetChatId()] = append(s.wheels[req.GetChatId()], req.GetOption())
    return &pb.AddOptionsResponse{Message: "Option added!"}, nil
}

func (s *server) SpinWheel(ctx context.Context, req *pb.SpinWheelRequest) (*pb.SpinWheelResponse, error) {
    options := s.wheels[req.GetChatId()]
    if len(options) == 0 {
        return &pb.SpinWheelResponse{Result: "No options!"}, nil
    }
    result := options[rand.Intn(len(options))]
    return &pb.SpinWheelResponse{Result: result}, nil
}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterFunWheelServiceServer(grpcServer, &server{wheels: make(map[int64][]string)})
    fmt.Printf("Server is running on port %s\n", port)
    grpcServer.Serve(lis)
}
