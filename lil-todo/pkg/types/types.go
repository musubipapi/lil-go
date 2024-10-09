package types

import "time"

type Todo struct {
	ID        int64     `json:"id"`
	Task      string    `json:"task"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}
