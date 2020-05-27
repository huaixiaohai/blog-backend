package dao

import (
	"database/sql"
	"fmt"
	"my/blog-backend/lib/log"
	"my/blog-backend/model"
	"time"
)

func createArticleInfoTable() {
	tableExist := "SELECT COUNT(*) FROM information_schema.TABLES WHERE table_name = ?;"

	var count int8
	var err error
	err = ConnDB.QueryRow(tableExist, tableArticleInfo).Scan(&count)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	if count != 0 {
		return
	}
	createArticleInfoSql := `CREATE TABLE blog_article_info (
		id BIGINT(20) NOT NULL AUTO_INCREMENT,
		title varchar(64) DEFAULT NULL,
		status INT(11) DEFAULT NULL,
		image LONGTEXT DEFAULT  NULL,
		summary TEXT DEFAULT NULL,
		author_id BIGINT(20) DEFAULT NULL ,
		content_id BIGINT(20) DEFAULT NULL,
		create_time DATETIME,
		update_time DATETIME,
		PRIMARY KEY (id)
	);`
	//var ret sql.Result
	_, err = ConnDB.Exec(createArticleInfoSql)
	if err != nil {
		log.Error(err)
		panic(err)
	}

}

type ArticleInfoDao struct {
}

var ArticleInfo = new(ArticleInfoDao)

func (a *ArticleInfoDao) Insert(title, summary, image, content string, authorID int64, status model.ArticleStatus) error {
	tx, err := ConnDB.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		err = tx.Rollback()
		if err != sql.ErrTxDone && err != nil {
			log.Error(err)
		}
	}()
	var result sql.Result
	result, err = tx.Exec("insert into blog_article_content (content, create_time, update_time) values (?, ?, ?)", content, time.Now(), time.Now())
	if err != nil {
		log.Error(err)
		return err
	}
	var contentID int64
	contentID, err = result.LastInsertId()
	if err != nil {
		log.Error(err)
		return err
	}
	_, err = tx.Exec("insert into blog_article_info (title, status, image, summary, author_id, content_id, create_time, update_time) values (?, ?, ?, ?, ?, ?, ?, ?)", title, status, image, summary, authorID, contentID, time.Now(), time.Now())
	if err != nil {
		log.Error(err)
		return err
	}
	if err = tx.Commit(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func (a *ArticleInfoDao) Update(ID int64, title, summary, image, content string, authorID int64, status model.ArticleStatus) error {
	tx, err := ConnDB.Begin()
	if err != nil {
		log.Error(err)
		return err
	}
	defer func() {
		err = tx.Rollback()
		if err != sql.ErrTxDone && err != nil {
			log.Error(err)
		}
	}()
	//var result sql.Result
	_, err = tx.Exec("update blog_article_info set title=?, status=?, image=?, summary=?, author_id=?, update_time=? where id=?", title, status, image, summary, authorID, time.Now(), ID)
	if err != nil {
		log.Error(err)
		return err
	}

	info := &model.ArticleInfo{}
	sqlStr := fmt.Sprintf("select * from %s where id = %d", tableArticleInfo, ID)
	tx.QueryRow(sqlStr).Scan(&info.ID, &info.Title, &info.Status, &info.Image, &info.Summary, &info.AuthorID, &info.ContentID, &info.CreateTime, &info.UpdateTime)

	_, err = tx.Exec("update blog_article_content set content=?, update_time=? where id=?", content, time.Now(), info.ContentID)
	if err != nil {
		log.Error(err)
		return err
	}
	if err = tx.Commit(); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

//func (a *ArticleInfoDao) Delete(title, summary, imgUrl, content string, authorID int64, status model.ArticleStatus) error {
//
//}

func (a *ArticleInfoDao) List(page, limit int64) (int64, []*model.ArticleInfo, error) {
	var err error
	sqlStr := fmt.Sprintf("select * from %s limit %d offset %d", tableArticleInfo, limit, limit*(page-1))
	var rows *sql.Rows
	rows, err = ConnDB.Query(sqlStr)
	defer func() {
		if rows != nil {
			rows.Close() //可以关闭掉未scan连接一直占用
		}
	}()
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	list := make([]*model.ArticleInfo, 0)
	for rows.Next() {
		info := &model.ArticleInfo{}
		err = rows.Scan(&info.ID, &info.Title, &info.Status, &info.Image, &info.Summary, &info.AuthorID, &info.ContentID, &info.CreateTime, &info.UpdateTime) //不scan会导致连接不释放
		if err != nil {
			log.Error(err)
			return 0, nil, err
		}
		list = append(list, info)
	}

	sqlStr = "select count(*) from " + tableArticleInfo
	var total int64 = 0
	err = ConnDB.QueryRow(sqlStr).Scan(&total)
	if err != nil {
		log.Error(err)
		return 0, nil, err
	}
	return total, list, nil
}

func (a *ArticleInfoDao) Detail(id int64) (*model.ArticleInfo, *model.ArticleContent, error) {
	sqlStr := fmt.Sprintf("select * from %s where id = %d", tableArticleInfo, id)
	info := &model.ArticleInfo{}
	err := ConnDB.QueryRow(sqlStr).Scan(&info.ID, &info.Title, &info.Status, &info.Image, &info.Summary, &info.AuthorID, &info.ContentID, &info.CreateTime, &info.UpdateTime)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	sqlStr = fmt.Sprintf("select * from %s where id = %d", tableArticleContent, info.ContentID)
	content := &model.ArticleContent{}
	err = ConnDB.QueryRow(sqlStr).Scan(&content.ID, &content.Content, &content.CreateTime, &content.UpdateTime)
	if err != nil {
		log.Error(err)
		return nil, nil, err
	}
	return info, content, nil
}
