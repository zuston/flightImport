package core

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"log"
	"bufio"
	"io"
	"strings"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"github.com/ivpusic/grpool"
)

const(
	// 型号
	MODEL = "model"
	// 架机
	AIR = "air"
	// 日期
	DATE = "date"
	// 架次
	SORTIE = "sortie"
	// 专业
	MAJOR = "major"
)

const(
	mysqlUserName = "root"
	mysqlPassword = "zuston"
)

var dataChan chan string
var endSemp chan bool

func init(){
	dataChan = make(chan string,100)
}

func DataSaver(path string, metaMapper map[string]string, sensorNames []string) {
	cf := "basic"
	var cFamilies = map[string]map[string]string{
		cf :  nil,
	}
	htablename := "pppp"
	// 判断hbase此表是否存在
	// 建表
	// 插入数据

	adminClient := gohbase.NewAdminClient("slave4,slave2,slave3")
	client := gohbase.NewClient("slave4,slave2,slave3")
	crt := hrpc.NewCreateTable(context.Background(), []byte(htablename), cFamilies)

	if err := adminClient.CreateTable(crt); err != nil {
		panic(err)
	}

	pool := grpool.NewPool(100, 100)
	// read from file
	defer pool.Release()
	go readLine(path)

	select {
	case v:=<-dataChan:
		columns := strings.Fields(v)
		rowkey := metaMapper[DATE]+"_"+columns[0]
		infoMapper := make(map[string][]byte,10)
		for i,v := range columns[1:]{
			sensorName := sensorNames[i+1]
			infoMapper[sensorName] = []byte(v)
		}
		basicInfoCfMapper := map[string]map[string][]byte{cf:infoMapper}

		pool.JobQueue <- func() {
			biPutRequest, err := hrpc.NewPutStr(context.Background(),htablename,rowkey,basicInfoCfMapper)
			if err!=nil {
				return
			}
			_, err = client.Put(biPutRequest)
			if err!=nil {
				return
			}
		}

	case <-endSemp:
		fmt.Println("finish")
		break
	}
}



func MetadataSaver(metaMapper map[string]string, sensorNames []string) {
	databaseName := "flight"
	model := metaMapper[MODEL]
	air := metaMapper[AIR]
	major := metaMapper[MAJOR]
	sortie := metaMapper[SORTIE]
	date := metaMapper[DATE]

	tableName := fmt.Sprintf("%s_%s_%s",model,air,major)
	MetaData := sensorNames

	// 判断是否存在此表
	// 建立表
	// 读取文件第一行，进行存储
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/?charset=utf8",mysqlUserName,mysqlPassword))
	defer db.Close()
	checkError(err)
	db.Query(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s;",databaseName))
	db.Query(fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s.%s (id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT,sortie varchar(20),date varchar(20),sensor varchar(20));",databaseName,tableName))

	for _,sensorName := range MetaData[1:]{
		// insert sql
		insertSql := fmt.Sprintf(`insert into %s.%s (sortie,date,sensor) values("%s","%s","%s");`,databaseName,tableName,sortie,date,sensorName)
		db.Query(insertSql)
	}
}

func readLine(path string){
	fi, err := os.Open(path)
	defer fi.Close()
	if err != nil {
		log.Panicf("Error: %s\n", err)
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			endSemp <- true
			break
		}
		dataChan <- string(a)
	}
}

func ReadMetaData(path string) []string{
	fi, err := os.Open(path)
	defer fi.Close()
	if err != nil {
		log.Panicf("Error: %s\n", err)
		return nil
	}
	defer fi.Close()

	br := bufio.NewReader(fi)
	for {
		a, _, c := br.ReadLine()
		if c == io.EOF {
			break
		}
		// 解析每行空格解析
		columns := strings.Fields(string(a))
		return columns
	}
	return nil
}

func checkError(e error) {
	if e!=nil {
		log.Panic(e)
	}
}


