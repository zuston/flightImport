package main

import (
	"github.com/zuston/flightImport/core"
	"fmt"
)



func main(){
	done := make(chan bool,1)
	path := "/opt/fuck.txt"
	// get filename and convert to mapper
	metaMapper := filenameParser(path)
	//
	sensorNames := core.ReadMetaData(path)
	// metadata to mysql
	core.MetadataSaver(metaMapper,sensorNames)
	// data to hbase
	core.DataSaver(path,metaMapper, sensorNames, done)

	<- done
	fmt.Println("处理完毕..........")
}

// 解析文件名
func filenameParser(path string) map[string]string {
	mapper := make(map[string]string)
	mapper["model"] = "kv"
	mapper["sortie"] = "1"
	mapper["date"] = "2018-09-09"
	mapper["air"] = "li"
	mapper["major"] = "test"
	return mapper
}

