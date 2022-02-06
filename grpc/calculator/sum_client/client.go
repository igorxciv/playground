package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/igorxciv/playground/grpc/calculator/sumpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Hello from sum client")

	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial to sum server: %v", err)
	}
	defer conn.Close()

	c := sumpb.NewSumServiceClient(conn)

	// doUnary(c)
	// doStream(c)

	doClientStream(c)
}

func doUnary(c sumpb.SumServiceClient) {
	fmt.Println("Starting unary RPC Sum request")
	req := &sumpb.SumRequest{
		FirstNumber:  5,
		SecondNumber: 391,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to sum: %v", err)
	}

	log.Printf("Response from sum: %v", res.Result)
}

func doStream(c sumpb.SumServiceClient) {
	fmt.Println("starting stream for prime point")
	req := &sumpb.PrimeRequest{
		Number: 120,
	}

	stream, err := c.Prime(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to get response:%v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to get data from response: %v", err)
		}
		log.Printf("response number: %v", msg.Number)
	}
}

func doClientStream(c sumpb.SumServiceClient) {
	log.Println("start client streaming")
	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("failed to do request as client stream")
	}

	reqs := []*sumpb.ComputeAverageRequest{
		{
			Number: 1,
		},
		{
			Number: 2,
		},
		{
			Number: 3,
		},
		{
			Number: 4,
		},
	}

	for _, req := range reqs {
		log.Printf("Sending number: %v", req.Number)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to close stream: %v", err)
	}

	log.Printf("AVERAGE: %v", res.Result)
}
