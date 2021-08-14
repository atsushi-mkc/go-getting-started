package models

import "time"

type User struct {
	ID        uint `gorm:"primary_key"`
	Name      string
	Age       int
	Height    *int
	Weight    *int
	UpdatedAt time.Time
	CreatedAt time.Time
}
