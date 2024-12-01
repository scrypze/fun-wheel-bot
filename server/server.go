package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "fun-wheel-bot/grpc" 
)

const port = ":50051"

type server struct {
	pb.UnimplementedFunWheelServiceServer
	wheels map[int64][]string
}

func (s *server) CreateWheel(ctx context.Context, req *pb.CreateWheelRequest) (*pb.CreateWheelResponse, error) {
	if _, exists := s.wheels[req.GetChatId()]; exists {
		return &pb.CreateWheelResponse{Message: "Колесо уже существует!"}, nil
	}
	s.wheels[req.GetChatId()] = []string{}
	return &pb.CreateWheelResponse{Message: "Колесо создано!"}, nil
}

func (s *server) AddOption(ctx context.Context, req *pb.AddOptionRequest) (*pb.AddOptionsResponse, error) {
	options, exists := s.wheels[req.GetChatId()]
	if !exists {
		return nil, fmt.Errorf("Колесо для этого чата не найдено")
	}
	s.wheels[req.GetChatId()] = append(options, req.GetOption())
	return &pb.AddOptionsResponse{Message: "Опция добавлена!"}, nil
}

func (s *server) SpinWheel(ctx context.Context, req *pb.SpinWheelRequest) (*pb.SpinWheelResponse, error) {
	options, exists := s.wheels[req.GetChatId()]
	if !exists || len(options) == 0 {
		return &pb.SpinWheelResponse{Result: "Колесо пустое!"}, nil
	}
	result := options[0] 
	return &pb.SpinWheelResponse{Result: result}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterFunWheelServiceServer(s, &server{wheels: make(map[int64][]string)})

	fmt.Println("Сервер запущен на порту :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
