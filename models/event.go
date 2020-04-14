package models

import (
	"time"
)

type Event struct {
	ID       string      `gorm:"type:text;unique_index" json:"id"`
	UserID   string      `gorm:"type:text;index" json:"user_id"`
	Date     time.Time   `json:"date"`
	Request  interface{} `gorm:"type:JSON;" json:"request"`
	Response interface{} `gorm:"type:JSON;" json:"response"`
	ApiType  APIType     `gorm:"type:text;" json:"api_type"`
}

type APIType int64

const (
	AScore = iota + 1
	BScore
)
