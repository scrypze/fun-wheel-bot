package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"math/rand"
	"time"

	"google.golang.org/grpc"
	pb "fun-wheel-bot/grpc"
)

type server struct {
	pb.UnimplementedFunWheelServiceServer
	wheels map[int64][]string
}

func (s *server) CreateWheel(ctx context.Context, req *pb.CreateWheelRequest) (*pb.CreateWheelResponse, error) {
	log.Printf("Creating wheel for chat_id: %d", req.GetChatId())
	s.wheels[req.GetChatId()] = make([]string, 0)
	return &pb.CreateWheelResponse{Message: "Колесо создано!"}, nil
}

func (s *server) AddOption(ctx context.Context, req *pb.AddOptionRequest) (*pb.AddOptionsResponse, error) {
	log.Printf("Adding option '%s' for chat_id: %d", req.GetOption(), req.GetChatId())
	
	if _, exists := s.wheels[req.GetChatId()]; !exists {
		return nil, fmt.Errorf("колесо не найдено, сначала создайте его")
	}
	
	s.wheels[req.GetChatId()] = append(s.wheels[req.GetChatId()], req.GetOption())
	return &pb.AddOptionsResponse{Message: "Опция добавлена!"}, nil
}

func (s *server) SpinWheel(ctx context.Context, req *pb.SpinWheelRequest) (*pb.SpinWheelResponse, error) {
	log.Printf("Spinning wheel for chat_id: %d", req.GetChatId())
	
	options, exists := s.wheels[req.GetChatId()]
	if !exists {
		return nil, fmt.Errorf("колесо не найдено, сначала создайте его")
	}
	
	if len(options) == 0 {
		return nil, fmt.Errorf("нет опций для выбора")
	}
	
	result := options[rand.Intn(len(options))]
	return &pb.SpinWheelResponse{Result: result}, nil
}

func main() {
	rand.Seed(time.Now().UnixNano())
	
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	
	s := grpc.NewServer()
	pb.RegisterFunWheelServiceServer(s, &server{
		wheels: make(map[int64][]string),
	})
	
	log.Printf("Server is running on port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
