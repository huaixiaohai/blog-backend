package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"time"
)

func main() {
	db, err := sql.Open("mysql", "root:2015.ami@tcp(localhost:3306)/test?charset=utf8")
	//数据库连接字符串，别告诉我看不懂。端口一定要写/
	if err != nil { //连接成功 err一定是nil否则就是报错
		panic(err.Error())       //抛出异常
		fmt.Println(err.Error()) //仅仅是显示异常
	}
	defer db.Close() //只有在前面用了 panic 这时defer才能起作用，如果链接数据的时候出问题，他会往err写数据

	rows, err := db.Query("select id, lvs from test")
	//判断err是否有错误的数据，有err数据就显示panic的数据
	if err != nil {
		panic(err.Error())
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()
	var id int        //定义一个id 变量
	var lvs string    //定义lvs 变量
	for rows.Next() { //开始循环
		rerr := rows.Scan(&id, &lvs) //数据指针，会把得到的数据，往刚才id 和 lvs引入
		if rerr == nil {
			fmt.Println("id号是", strconv.Itoa(id)+"     lvs lvs是"+lvs) //输出来而已，看看
		}
	}

	//_, e4 := db.Exec(insert_sql, "nima")
	go insert(db, "GO func")
	time.Sleep(1)
	insert(db, "Main Func")
	db.Close() //关闭数据库
}

func insert(db *sql.DB, s string) {
	insert_sql := "INSERT INTO test(lvs) VALUES(?)"

	for i := 1; i < 10; i++ {
		db.Exec(insert_sql, s)
	}
}