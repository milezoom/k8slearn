package main

import (
	"context"
	"fmt"
	"time"
	"net/http"
	_ "net/http/pprof"

	"printsvc/config"
	"printsvc/config/logger"
	"printsvc/config/repository"
	"printsvc/server"
	"printsvc/usecase"
	"printsvc/util"
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
	http.ListenAndServe("localhost" + restPort, nil)

	wait := util.GracefulShutdown(ctx, 5*time.Second, map[string]util.Operation{
		"grpc": func(ctx context.Context) error {
			grpcServer.GracefulStop()
			return nil
		},
	})
	<-wait
}
