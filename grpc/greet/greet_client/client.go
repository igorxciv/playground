package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/igorxciv/playground/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello from a client")

	certFile := "ssl/ca.crt"
	creds, _ := credentials.NewClientTLSFromFile(certFile, "")
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to dial to grpc: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	// doUnary(c)

	// doStream(c)

	// doClientStream(c)

	// doBidi(c)

	doWithDeadline(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Igor",
			LastName:  "Cheliadinski",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to get response: %v", err)
	}

	log.Printf("Response from greet: %v", res.Result)
}

func doStream(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do stream RPC")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Igor",
			LastName:  "Cheliadinski",
		},
	}

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to get response: %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("failed to get data from response: %v", err)
		}
		log.Printf("response from stream: %v", msg.GetResult())
	}
}

func doClientStream(c greetpb.GreetServiceClient) {
	fmt.Println("Streaming client started")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("failed to process client streaming: %v", err)
	}

	reqs := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Igor",
				LastName:  "Cheliadinski",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Victoria",
				LastName:  "Cheliadinskaya",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Frodo",
				LastName:  "Baggins",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Gendalf",
				LastName:  "Grey",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bilbo",
				LastName:  "Baggins",
			},
		},
	}

	for _, req := range reqs {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(1 * time.Second)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("failed to close and recive response: %v", err)
	}
	fmt.Printf("Long Greet response %v", res.Result)
}

func doBidi(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a BiDi Streaming RPC")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("failed to greet everyone: %v", err)
	}

	waitc := make(chan struct{})
	reqs := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Igor",
				LastName:  "Cheliadinski",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Victoria",
				LastName:  "Cheliadinskaya",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Frodo",
				LastName:  "Baggins",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Gendalf",
				LastName:  "Grey",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Bilbo",
				LastName:  "Baggins",
			},
		},
	}

	go func() {
		// send bunch of msgs
		for _, req := range reqs {
			fmt.Printf("Sending message: %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	go func() {
		// receive bunch of msgs
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("failed to receive response: %v", err)
			}
			fmt.Printf("Result: %v\n", res.Result)
		}
		close(waitc)
	}()

	<-waitc

}

func doWithDeadline(c greetpb.GreetServiceClient) {
	log.Println("Intercepted with deadline")

	req := &greetpb.GreetWithDeadlineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Igor",
			LastName:  "Cheliadinski",
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Println("Timeout!")
			} else {
				fmt.Printf("unexpected error: %v", statusErr.Message())
			}
		} else {
			log.Fatalf("Failed to make request: %v\n", err)
		}
		return
	}
	log.Printf("Response: %v\n", res.Result)
}
