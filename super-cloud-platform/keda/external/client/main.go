package main

import (
	"context"
	"log"
	"time"

	pb "github.com/jxs1211/external/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewExternalServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.Ping(ctx, &pb.PingRequest{Message: "Hello"})
	if err != nil {
		log.Fatalf("RPC failed: %v", err)
	}

	log.Printf("Reply: %s", res.Reply)
}
