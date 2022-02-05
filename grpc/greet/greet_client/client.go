package main

import (
	"fmt"
	"log"

	"github.com/igorxciv/playground/grpc/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Hello from a client")

	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to dial to grpc: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewDummyServiceClient(conn)
	fmt.Printf("created client: %f", c)
}
