package main

import (
	"context"
	"entgo.io/ent/dialect"
	_ "github.com/mattn/go-sqlite3"
	postv1 "github.com/sdoshi579/cloudbees/gen/post/v1"
	"github.com/sdoshi579/cloudbees/internal/repository/ent"
	postrepo "github.com/sdoshi579/cloudbees/internal/repository/post"
	postservice "github.com/sdoshi579/cloudbees/internal/service/post"
	postrpc "github.com/sdoshi579/cloudbees/rpc/post"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

func main() {
	logger := zap.NewExample()
	entClient, err := ent.Open(dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	if err := entClient.Schema.Create(context.Background()); err != nil {
		logger.Error("error in creating migration", zap.Error(err))
		os.Exit(1)
	}
	defer entClient.Close()

	logger.Info("initialized ent client")
	repository := postrepo.NewRepository(postrepo.WithEntClient(entClient), postrepo.WithLogger(logger))
	logger.Info("initialized repository")
	service := postservice.NewService(postservice.WithLogger(logger), postservice.WithRepository(repository))
	logger.Info("initialized service")
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen on port 50051: %v", err)
	}

	s := grpc.NewServer()

	postRPCInstance := postrpc.NewRPCImplementation(service, logger)
	postv1.RegisterPostServiceServer(s, postRPCInstance)

	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
