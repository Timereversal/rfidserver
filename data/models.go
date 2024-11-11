package data

import "database/sql"

type Models struct {
	Events EventModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Events: EventModel{DB: db},
	}
}
