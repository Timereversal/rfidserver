package main

import (
	"context"
	"github.com/Timereversal/rfidserver/reader"
	"google.golang.org/grpc"
	"log"
	"net"
)

type myReaderServer struct {
	reader.UnimplementedReaderServer
}

func (s myReaderServer) Report(ctx context.Context, request *reader.ReportRequest) (*reader.ReportResponse, error) {
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

	err = serverRegistrar.Serve(listen)
	if err != nil {
		log.Fatalf("error during server registrar %s", err)
	}
}
