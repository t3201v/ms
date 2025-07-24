package main

import (
	"net"
	"os"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	pb "github.com/t3201v/ms/resource/gen/resource/proto"
	"github.com/t3201v/ms/resource/internal/api"
	"github.com/t3201v/ms/resource/internal/log"
)

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", ":9090")

	lis, err := net.Listen("tcp", viper.GetString("PORT"))
	if err != nil {
		panic("failed to listen: " + err.Error())
	}

	if level := viper.GetString("LEVEL"); level == "prod" {
		log.Config(log.LevelDebug, log.OutputJSON, os.Stdout)
	}

	publicKey, err := api.LoadPublicKey()
	if err != nil {
		panic("failed to load public key: " + err.Error())
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			log.LoggingInterceptor,
			selector.UnaryServerInterceptor(api.AuthInterceptor(publicKey), selector.MatchFunc(api.SkipOnes)),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(func(p any) (err error) {
				log.Error(nil, "Recovered from panic", "reason", p)
				return status.Errorf(codes.Internal, "internal server error")
			})),
		),
	)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	pb.RegisterResourceServiceServer(s, api.NewResourceServer())

	reflection.Register(s)

	log.Info(nil, "Starting GRPC server", "port", viper.GetString("PORT"))
	if err := s.Serve(lis); err != nil {
		panic("failed to serve: " + err.Error())
	}
	s.GracefulStop()
}
