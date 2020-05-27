package model

import "time"

type ArticleStatus int8

const (
	ArticleStatusNone      ArticleStatus = 0
	ArticleStatusPublished ArticleStatus = 1 // 发布
	ArticleStatusDraft     ArticleStatus = 2 // 草案
)

type ArticleInfo struct {
	ID         int64
	Title      string        // 标题
	Status     ArticleStatus // 状态
	Image      string        // 图片
	Summary    string        // 摘要
	AuthorID   int64         // 作者
	ContentID  int64         // 内容
	CreateTime time.Time
	UpdateTime time.Time
}
