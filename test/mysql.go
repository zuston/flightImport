package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
)

//https://blog.csdn.net/jesseyoung/article/details/40398321

func main(){
	db, err := sql.Open("mysql", "root:zuston@tcp(127.0.0.1:3306)/?charset=utf8")
	if err!=nil {
		fmt.Println("connection error")
		return
	}
	db.Query("drop database if exists tmpdb")
	db.Query("create database tmpdb")
	//db.Query("use tmpdb")
	db.Query("create table tmpdb.tmptab(c1 int, c2 varchar(20), c3 varchar(20))")
	db.Query("insert into tmpdb.tmptab values(101, '姓名1', 'address1'), (102, '姓名2', 'address2'), (103, '姓名3', 'address3'), (104, '姓名4', 'address4')")
	//checkErr(err)
	query, err := db.Query("select * from tmpdb.tmptab")
	fmt.Print(err)
	for query.Next() {
		fmt.Println(query.Scan())
	}
	return
}
