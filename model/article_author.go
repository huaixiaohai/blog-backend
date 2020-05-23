package model

import "time"

type ArticleAuthor struct {
	ID         int64
	Name       string
	CreateTime *time.Time
	UpdateTime *time.Time
}
