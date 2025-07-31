package api

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "github.com/t3201v/ms/identity/gen/identity/proto"
	resource "github.com/t3201v/ms/identity/gen/resource/proto"
	"github.com/t3201v/ms/identity/internal/log"
)

type IdentityServer struct {
	pb.UnimplementedAuthServiceServer
	resourceUrl string
}

func (i *IdentityServer) Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := issueJWT("", nil)
	if err != nil {
		log.Error(ctx, "Error issuing JWT", "error", err)
		return nil, status.Errorf(codes.Internal, "failed to issue token: %v", err)
	}

	go func() {
		// test svc calls svc
		grpcClient, err := grpc.NewClient(i.resourceUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(err)
		}
		ctxWithToken := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer "+token)
		client := resource.NewResourceServiceClient(grpcClient)
		resp, err := client.SayHello2(ctxWithToken, &resource.HelloRequest{Name: "issuer"})
		if err != nil {
			panic(err)
		}
		log.Info(ctx, "Got response", "response", resp)
	}()

	return &pb.LoginResponse{
		Token: token,
	}, nil
}

func NewIdentityServer(resourceUrl string) pb.AuthServiceServer {
	return &IdentityServer{
		resourceUrl: resourceUrl,
	}
}
