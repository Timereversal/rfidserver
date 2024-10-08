package httpserver

import (
	"fmt"
	"github.com/Timereversal/rfidserver/pubsub"
	"net/http"
)

type SSEserver struct {
	Sub *pubsub.Server[pubsub.RunnerData]
}

func (s SSEserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Events Handler")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan pubsub.RunnerData)
	s.Sub.Subscribe(ch)
	defer s.Sub.Cancel(ch)

	for {
		select {
		case data := <-ch:
			newd := fmt.Sprintf(`data: {"runner":{"tagId":%d,"time_stage_1":"%s"}}`, data.TagId, data.TimeStage1)
			fmt.Fprintf(w, "%s\n\n", newd)
			w.(http.Flusher).Flush()

		}

	}

}
