package main

import (
	"context"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"hellosvc/config"
	"hellosvc/config/logger"
	"hellosvc/config/repository"
	"hellosvc/server"
	"hellosvc/usecase"
	"hellosvc/util"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
	config.LoadConfigMap()
	logger.LoadLogger()
	repository.LoadRepository()
}

func main() {
	// TODO: complete usecase implementation in usecase folder
	usecase := usecase.NewUsecase(repository.GetRepo())
	ctx := context.Background()
	grpcServer := server.RunGRPCServer(ctx, usecase)

	restPort := fmt.Sprintf(":%s", config.GetConfig("rest_port").GetString())
	http.ListenAndServe("localhost"+restPort, nil)

	wait := util.GracefulShutdown(ctx, 5*time.Second, map[string]util.Operation{
		"grpc": func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})
	<-wait
}
