package httpserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Timereversal/rfidserver/data"
	"github.com/Timereversal/rfidserver/pubsub"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"io"
	"net/http"
	"os"
)

type SSEserver struct {
	Sub *pubsub.Server[pubsub.RunnerData]
	DB  *sql.DB
}

func (s *SSEserver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Events Handler")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	ch := make(chan pubsub.RunnerData)
	s.Sub.Subscribe(ch)
	defer s.Sub.Cancel(ch)
	queryFormatString := `select re.tag_id,runner.name, runner.lastName, runner.distance, runner.sex, runner.category, runner.event_id,re.stage_0, re.stage_1, re.time_stage_1  
				from race_event_%d as re  
				left join runner on re.tag_id= runner.tag_id 
				where runner.event_id=%d and runner.tag_id=%d;`
	var Runner data.Runner

	for {
		select {
		case data := <-ch:
			query := fmt.Sprintf(queryFormatString, data.EventId, data.EventId, data.TagId)
			row := s.DB.QueryRow(query)
			err := row.Scan(&Runner.TagId, &Runner.Name, &Runner.LastName, &Runner.Distance, &Runner.Sex, &Runner.Category, &Runner.EventId, &Runner.Stage0, &Runner.Stage1, &Runner.TimeStage1)
			fmt.Println("")
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println(query)
			fmt.Printf("%+v\n", Runner)
			//b, err := json.Marshal(Runner)
			//fmt.Fprintln(w, `data: {"runner":`+string(b)+"}")
			//fmt.Fprintf(w, `data: {"runner": %s}\n\n`, string(b))
			newd := fmt.Sprintf(`data: {"runner":{"name":"%s"  ,"lastName":"%s","distance":%d,"sex":"%s","tagId":%d,"category":"%s","stage_0":"%s","stage_1":"%s","time_stage_1":"%s"}}`, Runner.Name, Runner.LastName, Runner.Distance, Runner.Sex, Runner.TagId, Runner.Category, Runner.Stage0, Runner.Stage1, Runner.TimeStage1)

			//newd := fmt.Sprintf(`data: {"runner":{"tagId":%d,"time_stage_1":"%s"}}`, data.TagId, data.TimeStage1)
			fmt.Fprintf(w, "%s\n\n", newd)
			w.(http.Flusher).Flush()

		}

	}

}

func UploadFile(w http.ResponseWriter, r *http.Request) {
	// Max file size 10MB

	fmt.Println("inside UploadFile")
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("runnerListFile")
	if err != nil {
		fmt.Printf("Error Retrieving the File: %v\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v, Size: %+v, \n", handler.Filename, handler.Size)

	dst, err := os.Create(handler.Filename)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func (s *SSEserver) CreateEventHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	//var input struct {
	//	EventName  string                   `json:"eventName"`
	//	Categories []map[string]interface{} `json:"categories,omitempty"`
	//	Distances  []int                    `json:"distances"`
	//}

	var event data.Event
	fmt.Printf("input: %+v\n", r.Body)
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		fmt.Printf("Error Decoding Input: %v\n", err)
		return
	}
	fmt.Println("event categories", event.Categories)
	jsonCategories, err := json.Marshal(event.Categories)
	fmt.Println(jsonCategories)
	if err != nil {
		fmt.Printf("Error Marshalling Categories: %v\n", err)
		return
	}
	query := `INSERT INTO events (name,time, categories, distance) VALUES ($1, $2, $3, $4) RETURNING id;`
	row := s.DB.QueryRow(query, event.EventName, event.EventDate, jsonCategories, event.Distances)
	var id int
	err = row.Scan(&id)
	if err != nil {
		fmt.Printf("Error Creating Event: %v\n", err)
	}

	createTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS race_event_%d (
		   tag_id INTEGER UNIQUE,
		   event_id INTEGER,
		   stage_0 TIMESTAMP WITH TIME ZONE,
		   stage_1 TIMESTAMP WITH TIME ZONE,
		   stage_2 TIMESTAMP WITH TIME ZONE,
		   stage_3 TIMESTAMP WITH TIME ZONE,
		   time_stage_1 INTERVAL GENERATED ALWAYS AS (AGE(stage_1,stage_0)) STORED,
		   time_stage_2 INTERVAL GENERATED ALWAYS AS (AGE(stage_2,stage_0)) STORED,
		   time_stage_3 INTERVAL GENERATED ALWAYS AS (AGE(stage_3,stage_0)) STORED
			);`, id)

	_, err = s.DB.Exec(createTable)
	if err != nil {
		fmt.Printf("Error Creating Table: %v\n", err)
	}
	fmt.Println("table created successfully")

	sqlStr := "INSERT INTO categories (name, age_low, age_high, fk_event_id) VALUES ($1, $2, $3, $4) ;"
	//vals := []interface{}{}
	tx, err := s.DB.Begin()
	if err != nil {
		return
	}

	for _, c := range event.Categories {
		_, err = tx.Exec(sqlStr, c["name"], c["ageLow"], c["ageHigh"], id)
		if err != nil {
			tx.Rollback()
			//fmt.Errorf("failed to insert user %s: %w", user.Name, err)
			fmt.Printf("failed to insert events")
			fmt.Println(err)
		}
		//vals = append(vals, c["name"], c["ageLow"], c["ageHigh"], id)
		//fmt.Println(v["name"], v["ageLow"], v["ageHigh"])
	}
	tx.Commit()
	//sqlStr = sqlStr[0 : len(sqlStr)-1]
	//stmt, _ := s.DB.Prepare(sqlStr)
	//res, _ := stmt.Exec(vals...)
	//fmt.Printf("Result %+v", res)
	//fmt.Printf("result %")
	fmt.Printf("input: %+v\n", event)
}

func (s *SSEserver) GetEventsInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")
	fmt.Printf("inside GetEventsInfo\n")
	var temp1 []byte

	var events []data.Event
	//query := `SELECT id, name, categories FROM events;`
	query := `SELECT id, name, distance, categories FROM events;`
	rows, err := s.DB.Query(query)
	//fmt.Println(rows)
	if err != nil {
		fmt.Printf("Error Retrieving Eventsss: %v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		//var numbers []int32
		var event data.Event
		//err = rows.Scan(&event.Id, &event.EventName, &event.CategoriesInfo)
		//err = rows.Scan(&event.Id, &event.EventName, pq.Array(&numbers), &temp1)
		//pq.Array support data type  https://github.com/lib/pq/blob/v1.10.9/array.go#L29
		err = rows.Scan(&event.Id, &event.EventName, pq.Array(&event.Distances), &temp1)
		if err != nil {
			fmt.Printf("Error Retrieving Events: %v\n", err)
		}
		err = json.Unmarshal(temp1, &event.Categories)

		events = append(events, event)
	}
	//fmt.Println(events)
	err = rows.Err()
	if err != nil {
		fmt.Printf("Error Retrieving Events: %v\n", err)
		return
	}
	c, _ := json.Marshal(events)
	w.Write(c)
}

func (s *SSEserver) CreateRunner(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	fmt.Println("inside CreateRunner")

	var runner data.Runner
	err := json.NewDecoder(r.Body).Decode(&runner)
	if err != nil {
		fmt.Printf("Error Decoding Input: %v\n", err)
	}
	fmt.Printf("input: %+v\n", runner)
	queryCategory := `SELECT * from categories WHERE fk_event_id = $1;`
	rows, err := s.DB.Query(queryCategory, runner.EventId)
	fmt.Println(rows)
	query := `INSERT INTO runner (name, lastname,birthDate,event_id, distance, sex, tag_id ) VALUES ($1, $2, $3, $4, $5,$6,$7) RETURNING id;`
	row := s.DB.QueryRow(query, runner.Name, runner.LastName, runner.BirthDate, runner.EventId, runner.Distance, runner.Sex, runner.TagId)
	var id int
	err = row.Scan(&id)
	if err != nil {
		fmt.Printf("Error Creating Runner: %v\n", err)
	}
}

func AssignTagId(eventId int32, category string) (int32, error) {
	return 0, nil
}

func (s *SSEserver) RunnersDataRace(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")

	var RunnersData []data.Runner

	queryFormatString := `select re.tag_id,runner.name, runner.lastName, runner.distance, runner.sex, runner.category, runner.event_id,re.stage_0, re.stage_1, re.time_stage_1  
				from race_event_37 as re  
				left join runner on re.tag_id= runner.tag_id 
				where runner.event_id=37;`

	rows, err := s.DB.Query(queryFormatString)
	if err != nil {
		fmt.Printf("Error Retrieving Events: %v\n", err)
	}
	defer rows.Close()

	for rows.Next() {
		var Runner data.Runner
		err = rows.Scan(&Runner.TagId, &Runner.Name, &Runner.LastName, &Runner.Distance, &Runner.Sex, &Runner.Category, &Runner.EventId, &Runner.Stage0, &Runner.Stage1, &Runner.TimeStage1)
		if err != nil {
			fmt.Printf("Error Retrieving Events: %v\n", err)
		}
		RunnersData = append(RunnersData, Runner)

	}
	err = rows.Err()
	if err != nil {
		fmt.Printf("Error Retrieving Events: %v\n", err)
		return
	}
	c, _ := json.Marshal(RunnersData)

	w.Write(c)

}
