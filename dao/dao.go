package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"my/blog-backend/conf"
	"my/blog-backend/lib/log"
	"time"
)

var (
	ConnDB *sql.DB
	err    error
)

const (
	tableUser = "blog_user"

	tableArticleAuthor  = "blog_article_author"
	tableArticleInfo    = "blog_article_info"
	tableArticleContent = "blog_article_content"
)

func init() {
	dbConfig := conf.C.DB
	ConnDB, err = sql.Open(conf.C.DB.Name, conf.C.DB.URL)
	if err != nil {
		panic(err.Error())
	}
	// 设置数据库最大连接数
	ConnDB.SetConnMaxLifetime(100 * time.Second)
	// 设置数据库最大闲置连接数
	ConnDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	// 验证连接
	err = ConnDB.Ping()
	if err != nil {
		log.Panic(err)
	}
	log.Info("数据库连接成功")

	createUserTable()
	createArticleAuthorTable()
	createArticleInfoTable()
	createArticleContentTable()
}
