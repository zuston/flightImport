package main

import "github.com/zuston/flightImport/core"



func main(){
	path := "./intro.md"
	// get filename and convert to mapper
	metaMapper := filenameParser(path)
	//
	sensorNames := core.ReadMetaData(path)
	// metadata to mysql
	core.MetadataSaver(metaMapper,sensorNames)
	// data to hbase
	core.DataSaver(path,metaMapper, sensorNames)

	select {

	}
}

// 解析文件名
func filenameParser(path string) map[string]string {
	return nil
}

