package model

import "time"

type Todo struct {
	ID        int
	Title     string
	DueDate   time.Time
	Completed bool
	Expired   bool
}
