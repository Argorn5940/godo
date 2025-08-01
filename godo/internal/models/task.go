package models

import "time"

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
	UpdatedAt   time.Time `json:"updated_at"`
}