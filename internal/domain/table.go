package domain

import (
	"encoding/json"
	"time"
)

type Table struct {
	Id          uint            `json:"id"`
	No          string          `json:"no"`
	GameId      uint            `json:"game_id"`
	Key         string          `json:"key"`
	State       uint8           `json:"state"`
	Description string          `json:"description"`
	Result      json.RawMessage `json:"result"`
	StartedAt   time.Time       `json:"started_at"`
	EndedAt     time.Time       `json:"ended_at"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	DeletedAt   time.Time       `json:"deleted_at"`
}
