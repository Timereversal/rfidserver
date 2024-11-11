package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/Timereversal/rfidserver/data"
	"github.com/Timereversal/rfidserver/httpserver"
	"github.com/Timereversal/rfidserver/pubsub"
	"github.com/Timereversal/rfidserver/reader"
	_ "github.com/jackc/pgx/v4/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

type application struct {
	models data.Models
}

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
	pubsub *pubsub.Server[pubsub.RunnerData]
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
	query := fmt.Sprintf("INSERT INTO race_event_%d(tag_id, %s) VALUES($1,$2)  ON CONFLICT(tag_id)  DO UPDATE SET %s = $2 RETURNING tag_id, time_stage_1 ;", request.EventId, stageId, stageId)
	fmt.Println(query)
	row := s.DbConn.QueryRow(query, request.TagId, runnerTime.Format("01-02-2006 15:04:05"))

	var id int
	var duration string
	err = row.Scan(&id, &duration)

	if err != nil {
		log.Println(err)
	}
	fmt.Println("tag_id", id, duration)

	s.pubsub.Publish(pubsub.RunnerData{TagId: id, TimeStage1: duration, StageId: int(request.Stage), EventId: request.EventId})

	return &reader.ReportResponse{
		Status: true,
	}, nil

}

type RFIDServer struct {
	httpHandler http.Handler
	grpcServer  myReaderServer
	eventStream chan pubsub.RunnerData
}

func main() {

	pubSubServer := pubsub.NewServer[pubsub.RunnerData]()

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

	//app := &application{
	//	models: data.NewModels(db),
	//}
	// http server initialization
	sseServ := &httpserver.SSEserver{Sub: pubSubServer, DB: db}
	//http.Handle()
	mux := http.NewServeMux()
	mux.Handle("/runners", sseServ)
	mux.HandleFunc("/upload", httpserver.UploadFile)
	mux.HandleFunc("/createEvent", sseServ.CreateEventHandler)
	mux.HandleFunc("/eventsInfo", sseServ.GetEventsInfo)
	mux.HandleFunc("/createRunner", sseServ.CreateRunner)
	mux.HandleFunc("/runnersDataRace", sseServ.RunnersDataRace)
	//mux.HandleFunc("/runners", eventsHandler)
	httpServer := &http.Server{
		Addr: net.JoinHostPort("0.0.0.0", "8090"),
		//Handler: app.routes(),
		Handler: mux,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Fprintf(os.Stderr, "error listening : %s\n", err)
		}
	}()
	//

	// Grpc start service
	listen, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("error during open tcp %s", err)
	}

	defer listen.Close()
	serverRegistrar := grpc.NewServer()
	service := &myReaderServer{DbConn: db, pubsub: pubSubServer}
	reader.RegisterReaderServer(serverRegistrar, service)
	reflection.Register(serverRegistrar)
	err = serverRegistrar.Serve(listen)
	if err != nil {
		log.Fatalf("error during server registrar %s", err)
	}

	//http.HandleFunc("/runners", eventsHandler)
	//http.ListenAndServe(":8090", nil)

}
