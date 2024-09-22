package main

import (
	"context"
	"fmt"
	"github.com/Timereversal/rfidserver/reader"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
)

type myReaderServer struct {
	//reader.UnimplementedReaderServer
	reader.UnimplementedReaderServer
}

func (s myReaderServer) Report(ctx context.Context, request *reader.ReportRequest) (*reader.ReportResponse, error) {
	fmt.Printf("TagId: %d , EventID: %d, time:%s \n", request.TagId, request.EventId, request.RunnerTime.AsTime().Format("01-02-2006 15:04:05"))
	return &reader.ReportResponse{
		Status: true,
	}, nil
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Events Handler")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	for i := 0; i < 20; i++ {
		//fmt.Fprintf(w, "data: %d\n\n", i)
		newd := fmt.Sprintf(`data: {"runner":{"tagId":%d,"eventId":%d}}`, i, i)
		fmt.Fprintf(w, "%s\n\n", newd)
		w.(http.Flusher).Flush()
		time.Sleep(2 * time.Second)
	}
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
	reflection.Register(serverRegistrar)
	err = serverRegistrar.Serve(listen)
	if err != nil {
		log.Fatalf("error during server registrar %s", err)
	}

	//http.HandleFunc("/runners", eventsHandler)
	//http.ListenAndServe(":8090", nil)

}
