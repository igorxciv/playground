package main

import (
	"context"
	"fmt"
	"log"

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

	doUnary(c)
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
