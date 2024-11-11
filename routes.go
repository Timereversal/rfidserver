package main

import (
	"github.com/Timereversal/rfidserver/httpserver"
	"net/http"
)

func (app *application) routes() *http.ServeMux {

	//sseServ := httpserver.SSEserver{Sub: pubsubserver}
	mux := http.NewServeMux()

	//mux.Handle("/runners", sseServ)
	mux.HandleFunc("/upload", httpserver.UploadFile)
	//mux.HandleFunc("/createEvent", .CreateEventHandler)
	return mux
}
