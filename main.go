package main

import (
	"context"
	"fmt"
	"github.com/Timereversal/rfidserver/reader"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type myReaderServer struct {
	reader.UnimplementedReaderServer
}

func (s myReaderServer) Report(ctx context.Context, request *reader.ReportRequest) (*reader.ReportResponse, error) {
	fmt.Printf("Runner Id: %d, EventId: %d", request.GetTagId(), request.GetEventId())
	return &reader.ReportResponse{
		Status: true,
	}, nil
}

func main() {
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("error during open tcp %s", err)
	}

	defer listen.Close()

	serverRegistrar := grpc.NewServer()
	service := &myReaderServer{}

	reader.RegisterReaderServer(serverRegistrar, service)
	// enable reflection GRPC
	reflection.Register(serverRegistrar)
	err = serverRegistrar.Serve(listen)
	if err != nil {
		log.Fatalf("error during server registrar %s", err)
	}
}
