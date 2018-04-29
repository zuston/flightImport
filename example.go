package main

import (
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"context"
)

func main(){
	var cFamilies = map[string]map[string]string{
		"basic" :  nil,
	}
	adminClient := gohbase.NewAdminClient("slave4,slave2,slave3")
	client := gohbase.NewClient("slave4,slave2,slave3")
	crt := hrpc.NewCreateTable(context.Background(), []byte("htest"), cFamilies)

	if err := adminClient.CreateTable(crt); err != nil {
		panic(err)
	}

	basicInfoCfMapper := map[string]map[string][]byte{"basic":{"name":[]byte("error")}}

	biPutRequest, err := hrpc.NewPutStr(context.Background(),"htest","123",basicInfoCfMapper)
		if err!=nil {
			return
		}
		_, err = client.Put(biPutRequest)
		if err!=nil {
			return
		}
}
