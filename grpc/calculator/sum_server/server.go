package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

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

func (s *server) Prime(req *sumpb.PrimeRequest, stream sumpb.SumService_PrimeServer) error {
	k := int32(2)
	v := req.GetNumber()
	for v > 1 {
		if v%k == 0 {
			stream.Send(&sumpb.PrimeResponse{
				Number: k,
			})
			time.Sleep(1 * time.Second)
			v = v / 2
		} else {
			k = k + 1
		}
	}
	return nil
}

func (s *server) ComputeAverage(stream sumpb.SumService_ComputeAverageServer) error {
	log.Println("started computing average...")

	sum := 0
	l := 0
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&sumpb.ComputeAverageResponse{
				Result: float32(sum) / float32(l),
			})
		}
		if err != nil {
			log.Fatalf("failed to receive from client stream: %v", err)
		}
		sum += int(msg.GetNumber())
		l = l + 1
	}
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
