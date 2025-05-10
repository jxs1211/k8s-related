package main

import (
	"context"
	"log"
	"net"

	"github.com/jxs1211/external/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.UnimplementedExternalScalerServer
}

func (c *server) IsActive(ctx context.Context, in *pb.ScaledObjectRef) (*pb.IsActiveResponse, error) {
	log.Printf("Received: %v", in)
	return &pb.IsActiveResponse{Result: true}, nil
}

func (c *server) GetMetricSpec(ctx context.Context, in *pb.ScaledObjectRef) (*pb.GetMetricSpecResponse, error) {
	log.Printf("Received: %v", in)
	return &pb.GetMetricSpecResponse{MetricSpecs: []*pb.MetricSpec{
		{
			MetricName: "test",
			TargetSize: 100,
		},
	}}, nil
}

func (c *server) GetMetrics(ctx context.Context, in *pb.GetMetricsRequest) (*pb.GetMetricsResponse, error) {
	log.Printf("Received: %v", in)
	return &pb.GetMetricsResponse{MetricValues: []*pb.MetricValue{
		{
			MetricName:  "test",
			MetricValue: 100,
		},
	}}, nil
}

func (c *server) StreamIsActive(in *pb.ScaledObjectRef, srv pb.ExternalScaler_StreamIsActiveServer) error {
	log.Printf("Received: %v", in)
	return status.Errorf(codes.Unimplemented, "method StreamIsActive not implemented")
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterExternalScalerServer(s, &server{})

	log.Println("Server running on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
