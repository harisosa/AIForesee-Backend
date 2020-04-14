package models

import "time"

type User struct {
	ID       string    `gorm:"primary_key;unique_index,not_null" json:"id"`
	Username string    `gorm:"type:varchar(255);NOT NULL;unique_index" json:"username"`
	Password string    `gorm:"type:varchar(255);NOT NULL" json:"password"`
	Email    string    `gorm:"type:varchar(255);NOT NULL;unique_index" json:"email"`
	ApiKey   string    `gorm:"type:varchar(255);NOT NULL;" json:"api_key"`
	Created  time.Time `json:"created"`
}
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
