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
	"google.golang.org/grpc/status"
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

	// doClientStream(c)

	// doBidi(c)

	doErrorUnary(c)
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

func doBidi(c sumpb.SumServiceClient) {
	fmt.Println("Client bidi started!")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("failed to get maximum: %v", err)
	}

	waitc := make(chan struct{})
	reqs := []*sumpb.FindMaximumRequest{
		{
			Number: 1,
		},
		{
			Number: 5,
		},
		{
			Number: 3,
		},
		{
			Number: 6,
		},
		{
			Number: 2,
		},
		{
			Number: 2,
		},
		{
			Number: 20,
		},
	}

	go func() {
		for _, req := range reqs {
			fmt.Printf("Sending to get max: %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("failed to receive response : %v", err)
			}
			fmt.Printf("Maximum for now: %v\n", res.Result)
		}
		close(waitc)
	}()

	<-waitc
}

func doErrorUnary(c sumpb.SumServiceClient) {
	fmt.Println("Starting square root unary RPC...")

	res, err := c.SquareRoot(context.Background(), &sumpb.SquareRootRequest{
		Number: -2,
	})
	if err != nil {
		resErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			log.Println(resErr.Message())
		} else {
			log.Fatalf("failed calling SquareRoot: %v", err)
		}
	}

	log.Printf("Square: %v", res.GetNumber())
}
