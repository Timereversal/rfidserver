package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Timereversal/rfidserver/reader"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"time"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

type myReaderServer struct {
	//reader.UnimplementedReaderServer
	reader.UnimplementedReaderServer
	DbConn *sql.DB
}

func (s myReaderServer) Report(ctx context.Context, request *reader.ReportRequest) (*reader.ReportResponse, error) {

	lima, err := time.LoadLocation("America/Lima")
	if err != nil {
		fmt.Println(err)
	}
	runnerTime := request.RunnerTime.AsTime().In(lima)
	fmt.Printf("TagId: %d ,Stage: %d, EventID: %d, time:%s \n", request.TagId, request.Stage, request.EventId, runnerTime.Format("01-02-2006 15:04:05"))

	stageId := fmt.Sprintf("stage_%d", request.Stage)
	//columnNameStage := fmt.Sprintf("_%d",request.Stage)
	query := fmt.Sprintf("INSERT INTO race_2345(tag_id, %s) VALUES($1,$2) ON CONFLICT(tag_id) DO UPDATE SET %s = $2;", stageId, stageId)
	_, err = s.DbConn.Exec(query, request.TagId, runnerTime.Format("01-02-2006 15:04:05"))
	if err != nil {
		log.Println(err)
	}
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

	postgresCfg := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "baloo",
		Password: "junglebook",
		Database: "lenslocked",
		SSLMode:  "disable",
	}

	db, err := sql.Open("pgx", postgresCfg.String())
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Database connection established")

	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("error during open tcp %s", err)
	}

	defer listen.Close()
	serverRegistrar := grpc.NewServer()
	service := &myReaderServer{DbConn: db}
	reader.RegisterReaderServer(serverRegistrar, service)
	reflection.Register(serverRegistrar)
	err = serverRegistrar.Serve(listen)
	if err != nil {
		log.Fatalf("error during server registrar %s", err)
	}

	//http.HandleFunc("/runners", eventsHandler)
	//http.ListenAndServe(":8090", nil)

}
