package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Name      string    `json:"name"`
	Age       int       `json:"age"`
	Height    *int      `json:"height"`
	Weight    *int      `json:"weight"`
	UpdatedAt time.Time `json:"-"`
	CreatedAt time.Time `json:"-"`
}
