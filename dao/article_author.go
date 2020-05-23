package dao

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"my/blog-backend/lib/log"
	"my/blog-backend/model"
)

func createArticleAuthorTable() {
	tableExist := "SELECT COUNT(*) FROM information_schema.TABLES WHERE table_name = ?;"

	var count int8
	var err error
	err = ConnDB.QueryRow(tableExist, tableArticleAuthor).Scan(&count)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	if count != 0 {
		return
	}
	createArticleAuthorSql := `CREATE TABLE blog_article_author (
		id BIGINT(20) NOT NULL AUTO_INCREMENT,
		name varchar(64) DEFAULT NULL,
		create_time DATETIME,
		update_time DATETIME,
		PRIMARY KEY (id)
	);`
	//var ret sql.Result
	_, err = ConnDB.Exec(createArticleAuthorSql)
	if err != nil {
		log.Error(err)
		panic(err)
	}

}

type ArticleAuthorDao struct {
}

var ArticleAuthor = new(ArticleAuthorDao)

func (u *ArticleAuthorDao) List() ([]*model.ArticleAuthor, error) {
	var err error

	sqlStr := "select * from " + tableArticleAuthor
	var rows *sql.Rows
	rows, err = ConnDB.Query(sqlStr)
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		log.Error(err)
		return nil, err
	}
	authors := make([]*model.ArticleAuthor, 0)
	for rows.Next() {
		author := &model.ArticleAuthor{}
		err = rows.Scan(&author.ID, &author.Name, &author.CreateTime, &author.UpdateTime) //不scan会导致连接不释放
		if err != nil {
			log.Error(err)
			return nil, err
		}
		authors = append(authors, author)
	}
	return authors, nil
}
