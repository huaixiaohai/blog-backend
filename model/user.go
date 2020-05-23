package model

import "time"

type User struct {
	ID         int64
	UserName   string
	Password   string
	CreateTime *time.Time
	UpdateTime *time.Time
}
