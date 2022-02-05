package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/igorxciv/playground/grpc/calculator/sumpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Sum(ctx context.Context, req *sumpb.SumRequest) (*sumpb.SumResponse, error) {
	log.Printf("sum server is called with %v, %v", req.FirstNumber, req.SecondNumber)
	res := &sumpb.SumResponse{
		Result: req.FirstNumber + req.SecondNumber,
	}
	return res, nil
}

func main() {
	fmt.Println("Hello from sum server")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen on port 50051")
	}

	s := grpc.NewServer()
	sumpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
