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

const port = ":50051"

type server struct {
    pb.UnimplementedFunWheelServiceServer
    wheels map[int64][]string
    rnd    *rand.Rand
}

func (s *server) CreateWheel(ctx context.Context, req *pb.CreateWheelRequest) (*pb.CreateWheelResponse, error) {
    log.Printf("Received CreateWheel request with chat_id: %d", req.GetChatId())
    if s.wheels == nil {
        log.Printf("Initializing wheels map")
        s.wheels = make(map[int64][]string)
    }
    s.wheels[req.GetChatId()] = make([]string, 0)
    log.Printf("Wheel created for chat_id: %d", req.GetChatId())
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
    
    result := options[s.rnd.Intn(len(options))]
    return &pb.SpinWheelResponse{Result: result}, nil
}

func (s *server) ViewOptions(ctx context.Context, req *pb.ViewOptionsRequest) (*pb.ViewOptionsResponse, error) {
    log.Printf("Viewing options for chat_id: %d", req.GetChatId())
    
    options, exists := s.wheels[req.GetChatId()]
    if !exists {
        return nil, fmt.Errorf("колесо не найдено, сначала создайте его")
    }
    
    return &pb.ViewOptionsResponse{Options: options}, nil
}

func main() {
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }
    
    s := grpc.NewServer()
    
    pb.RegisterFunWheelServiceServer(s, &server{
        wheels: make(map[int64][]string),
        rnd:    r,
    })
    
    log.Printf("Server is running on port %s", port)
    
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
