package main

import (
	"net"

	"github.com/spf13/viper"
	pb "github.com/t3201v/ms/identity/gen/identity/proto"
	"github.com/t3201v/ms/identity/internal/api"
	"github.com/t3201v/ms/identity/internal/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", ":9090")
	viper.SetDefault("RESOURCE_URL", "resource-service:9090")

	lis, err := net.Listen("tcp", viper.GetString("PORT"))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	s := grpc.NewServer(
		grpc.UnaryInterceptor(log.LoggingInterceptor),
	)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	pb.RegisterAuthServiceServer(s, api.NewIdentityServer(viper.GetString("RESOURCE_URL")))

	reflection.Register(s)

	log.Info(nil, "Starting GRPC server", "port", viper.GetString("PORT"))
	if err := s.Serve(lis); err != nil {
		panic("failed to serve: " + err.Error())
	}
	s.GracefulStop()
}
