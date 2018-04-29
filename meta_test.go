package main

import (
	"testing"
	"github.com/zuston/flightImport/core"
)

func TestMetaSaver(t *testing.T) {
	mapper := make(map[string]string)
	mapper["model"] = "kv"
	mapper["sortie"] = "1"
	mapper["date"] = "2018-09-09"
	mapper["air"] = "li"
	mapper["major"] = "test"

	path := "./data/fuck.txt"
	core.MetadataSaver(path,mapper)
}
