package main

import (
	"context"
	"log"
	"time"

	pb "github.com/jxs1211/external/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const grpcConfig = `{"loadBalancingConfig": [{"round_robin":{}}]}`

func main() {
	// addr := "localhost:50051"
	addr := "external-scaler.default.svc.cluster.local:50051"
	conn, err := grpc.NewClient(addr,
		grpc.WithDefaultServiceConfig(grpcConfig),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewExternalScalerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := c.IsActive(ctx, &pb.ScaledObjectRef{})
	if err != nil {
		log.Fatalf("RPC failed: %v", err)
	}

	log.Printf("Reply: %s", res.Result)
	msResp, err := c.GetMetricSpec(ctx, &pb.ScaledObjectRef{
		Name:           "test",
		Namespace:      "test",
		ScalerMetadata: map[string]string{"key": "val"},
	})
	if err != nil {
		log.Fatalf("RPC failed: %v", err)
	}

	log.Printf("Reply: %v", msResp)
}
