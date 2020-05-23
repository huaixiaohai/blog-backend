package model

import "time"

type ArticleContent struct {
	ID         int64
	Content    string
	CreateTime time.Time
	UpdateTime time.Time
}
