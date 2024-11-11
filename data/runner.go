package data

import "time"

type Runner struct {
	Name         string    `json:"name"`
	LastName     string    `json:"lastName"`
	BirthDate    time.Time `json:"birthDate"`
	Distance     int32     `json:"distance"`
	Sex          string    `json:"sex"`
	Category     string    `json:"category"`
	TagId        int32     `json:"tagId"`
	EventId      int32     `json:"eventId"`
	RankCategory int32     `json:"rankCategory"`
	RankAll      int32     `json:"rankAll"`
	Stage0       time.Time `json:"stage_0"`
	Stage1       time.Time `json:"stage_1"`
	TimeStage1   string    `json:"time_stage_1"`
}
