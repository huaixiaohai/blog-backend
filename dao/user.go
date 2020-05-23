package dao

import (
	"my/blog-backend/lib/log"
	"my/blog-backend/model"
)

func createUserTable() {
	tableExist := "SELECT COUNT(*) FROM information_schema.TABLES WHERE table_name = ?;"

	var count int8
	var err error
	err = ConnDB.QueryRow(tableExist, tableUser).Scan(&count)
	if err != nil {
		log.Error(err)
		panic(err)
	}
	if count != 0 {
		return
	}
	createUserSql := `CREATE TABLE blog_user (
		id BIGINT(20) NOT NULL AUTO_INCREMENT,
		user_name varchar(64) DEFAULT NULL,
		password  varchar(64) DEFAULT NULL,
		create_time DATETIME,
		update_time DATETIME,
		PRIMARY KEY (id)
	);`
	//var ret sql.Result
	_, err = ConnDB.Exec(createUserSql)
	if err != nil {
		log.Error(err)
		panic(err)
	}

}

type UserDao struct {
}

var User = new(UserDao)

func (u *UserDao) One(userName string) (*model.User, error) {
	var err error
	user := &model.User{}
	sql := "select * from " + tableUser + " where user_name = '" + userName + "'"
	err = ConnDB.QueryRow(sql).Scan(&user.ID, &user.UserName, &user.Password, &user.CreateTime, &user.UpdateTime)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return user, nil
}
