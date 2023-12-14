package server

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"runtime/debug"

	"hellosvc/config"
	"hellosvc/config/logger"
	"hellosvc/contract"
	"hellosvc/transport"
	"hellosvc/usecase"

	commonLogger "git.bluebird.id/mybb-ms/aphrodite/logger"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.elastic.co/apm/module/apmgrpc/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func RunGRPCServer(ctx context.Context, usecase *usecase.UseCase) *grpc.Server {
	grpcPort := fmt.Sprintf(":%s", config.GetConfig("grpc_port").GetString())
	grpcConn, err := net.Listen("tcp", grpcPort)
	if err != nil {
		logger.GetLogger().Fatal(fmt.Sprintf("failed to listen port: %v", err))
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			apmgrpc.NewUnaryServerInterceptor(),
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(grpcRecoveryHandler)),
		),
		grpc.ChainStreamInterceptor(
			apmgrpc.NewStreamServerInterceptor(),
			recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(grpcRecoveryHandler)),
		),
	)
	transport := transport.NewTransport(ctx, usecase)
	contract.RegisterHelloServiceServer(grpcServer, transport)
	reflection.Register(grpcServer)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)
	for name := range grpcServer.GetServiceInfo() {
		healthServer.SetServingStatus(name, healthpb.HealthCheckResponse_SERVING)
	}
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go grpcServer.Serve(grpcConn)
	logger.GetLogger().Info(fmt.Sprintf("server grpc listening at %v", grpcConn.Addr()))
	return grpcServer
}

func grpcRecoveryHandler(ctx context.Context, panic interface{}) error {
	newLineRegex := regexp.MustCompile(`\r?\n`)
	stackTrace := newLineRegex.ReplaceAllString(string(debug.Stack()), " ")
	logger.GetLogger().Error("panic happened", commonLogger.ConvertMapToFields(map[string]interface{}{
		"panic_message":    panic,
		"panic_stacktrace": stackTrace,
	})...)
	return status.Errorf(codes.Internal, "server error happened")
}
