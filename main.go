package main

import (
	"blockpost/config"
	"blockpost/genprotos/authorization"
	services "blockpost/services/authorization"
	"blockpost/storage"
	"blockpost/storage/postgres"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	conf := config.Load()
	AUTH := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.PostgresHost,
		conf.PostgresPort,
		conf.PostgresUser,
		conf.PostgresPassword,
		conf.PostgresDatabase,
	)
	var inter storage.StorageI
	inter, err := postgres.InitDB(AUTH)
	if err != nil {
		panic(err)
	}

	fmt.Printf("gRPC server running port%s with tcp protocol!", conf.GRPCPort)

	listener, err := net.Listen("tcp", conf.GRPCPort)
	if err != nil {
		panic(err)
	}

	srv := grpc.NewServer()

	authService := services.NewAuthService(conf, inter)
	authorization.RegisterAuthServiceServer(srv, authService)

	reflection.Register(srv)

	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
