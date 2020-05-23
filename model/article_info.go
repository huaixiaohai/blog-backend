package model

import "time"

type ArticleStatus int8

const (
	ArticleStatusNone      ArticleStatus = 0
	ArticleStatusPublished ArticleStatus = 1 // 发布
	ArticleStatusDraft     ArticleStatus = 2 // 草案
)

type ArticleInfo struct {
	ID            int64
	Title         string        // 标题
	Status        ArticleStatus // 状态
	ImgUrl        string        // 图片url
	Summary       string        // 摘要
	AuthorID      int64         // 作者
	ContentID     int64         // 内容
	PublishedTime time.Time     // 发布时间
	CreateTime    time.Time
	UpdateTime    time.Time
}
