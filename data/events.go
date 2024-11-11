package data

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type Event struct {
	Id         int
	EventName  string                   `json:"eventName"`
	EventDate  time.Time                `json:"eventDate"`
	Categories []map[string]interface{} `json:"categories,omitempty"`
	Distances  []int32                  `json:"distances"`
}

type EventInfo struct {
	Id             int        `json:"id" db:"id"`
	EventName      string     `json:"eventName" db:"name"`
	CategoriesInfo Categories `json:"categories" db:"categories"`
	//CategoriesInfo []CategoryItem
	Distances []int `json:"distances" db:"distance"`
}

type Categories struct {
	Categories []CategoryItem `json:"categories"`
}
type CategoryItem struct {
	Name         string `json:"name"`
	AgeLow       int    `json:"ageLow"`
	AgeHigh      int    `json:"ageHigh"`
	Participants int    `json:"participants"`
}

type EventModel struct {
	DB *sql.DB
}

func (c *Categories) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Categories) Scan(value interface{}) error {
	var temp Categories
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	fmt.Println(b)
	return json.Unmarshal(b, &temp)
}

func (e EventModel) Create(event *Event) error {

	jsonCategories, err := json.Marshal(event.Categories)
	if err != nil {
		return err
	}
	query := `INSERT INTO events (name, categories, distance) VALUES ($1, $2, $3) RETURNING id;`
	row := e.DB.QueryRow(query, event.EventName, jsonCategories, event.Distances)
	var id int
	err = row.Scan(&id)
	if err != nil {
		return err
	}
	return nil
}
func (e EventModel) Get(id int64) (*Event, error) {
	return nil, nil
}

func (e EventModel) Update(event *Event) error {
	return nil
}

func (e EventModel) Delete(id int64) error {
	return nil
}

func (e EventModel) getAllEvents() ([]Event, error) {

	var events []Event
	query := `SELECT id, name, categories, distance FROM events;`
	rows, err := e.DB.Query(query)
	if err != nil {
		fmt.Println("Error getting all events: ", err)
		return nil, nil
	}
	defer rows.Close()
	for rows.Next() {
		var e Event
		err := rows.Scan(&e.Id, &e.EventName, &e.Categories, &e.Distances)
		if err != nil {
			fmt.Println("Error getting all events: ", err)
			return nil, nil
		}
		events = append(events, e)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	return nil, nil
}
