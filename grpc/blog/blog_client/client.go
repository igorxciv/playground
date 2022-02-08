package main

import (
	"context"
	"log"

	"github.com/igorxciv/playground/grpc/blog/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	log.Println("Starting BLOG Client...")

	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	cc, err := grpc.Dial("0.0.0.0:50051", opts)
	if err != nil {
		log.Fatalf("failed to connect to gRPC server: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// create
	log.Println("Creating the blog post...")
	req := &blogpb.CreateBlogRequest{
		Blog: &blogpb.Blog{
			AuthorId: "Igor",
			Title:    "This is blog created from gRPC",
			Content:  "gRPC is awesome!",
		},
	}

	res, err := c.CreateBlog(context.Background(), req)
	if err != nil {
		log.Fatalf("failed to create a blog post: %v", err)
	}

	log.Printf("Blog has been created: %v\n", res.GetBlog())

	// read
	log.Println("Reading the blog post...")
	readReq := &blogpb.ReadBlogRequest{
		BlogId: res.GetBlog().GetId(),
	}
	readRes, err := c.ReadBlog(context.Background(), readReq)
	if err != nil {
		log.Fatalf("failed to read a blog post: %v", err)
	}

	log.Printf("Successfully read a blog post: %v\n", readRes.GetBlog())
}
