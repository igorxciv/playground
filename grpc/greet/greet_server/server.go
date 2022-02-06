package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/igorxciv/playground/grpc/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct{}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("greet function was invoked with %v", req)
	fName := req.GetGreeting().GetFirstName()
	result := "Hello " + fName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + " number" + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(2 * time.Second)
	}
	return nil
}

func (s *server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}
		if err != nil {
			log.Fatalf("failed to get request: %v", err)
		}
		fName := req.GetGreeting().GetFirstName()
		lName := req.GetGreeting().GetLastName()
		result += "Hello " + fName + " " + lName + "!\n"
	}
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
