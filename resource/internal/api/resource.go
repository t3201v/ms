package api

import (
	"context"

	pb "github.com/t3201v/ms/resource/gen/resource/proto"
	"github.com/t3201v/ms/resource/internal/log"
)

type ResourceServer struct {
	pb.UnimplementedResourceServiceServer
}

func (r *ResourceServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Info(ctx, "Received", "in", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func (r *ResourceServer) SayHello2(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Info(ctx, "Received v2", "in", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func NewResourceServer() pb.ResourceServiceServer {
	return &ResourceServer{}
}
