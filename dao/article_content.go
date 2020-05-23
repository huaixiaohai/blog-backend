package dao

import (
	"my/blog-backend/lib/log"
)

func createArticleContentTable() {
	tableExist := "SELECT COUNT(*) FROM information_schema.TABLES WHERE table_name = ?;"

	var count int8
	var err error
	err = ConnDB.QueryRow(tableExist, tableArticleContent).Scan(&count)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	if count != 0 {
		return
	}
	createArticleContentSql := `CREATE TABLE blog_article_content (
		id BIGINT(20) NOT NULL AUTO_INCREMENT,
		content TEXT DEFAULT NULL ,
		create_time DATETIME,
		update_time DATETIME,
		PRIMARY KEY (id)
	);`
	//var ret sql.Result
	_, err = ConnDB.Exec(createArticleContentSql)
	if err != nil {
		log.Error(err)
		panic(err)
	}

}

type ArticleContentDao struct {
}

var ArticleContent = new(ArticleContentDao)

func (a *ArticleContentDao) InsertArticleContent() {
	
}
