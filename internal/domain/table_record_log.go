package domain

import (
	"encoding/json"
	"time"
)

type TableRecordLog struct {
	Id            uint            `json:"id"`
	GameId        uint            `json:"game_id"`
	TableRecordId uint            `json:"table_record_id"`
	State         uint8           `json:"state"`
	Text          json.RawMessage `json:"text"`
	Image         json.RawMessage `json:"image"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
	DeletedAt     time.Time       `json:"deleted_at"`
}
