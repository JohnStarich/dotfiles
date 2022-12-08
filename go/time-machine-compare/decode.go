package main

import (
	"encoding/json"
	"encoding/xml"
	"io"

	"github.com/johnstarich/go/plist"
)

type TMUtilCompare struct {
	XMLName xml.Name `xml:"dict"`
	Changes []Change
	Totals  Totals
}

type Change struct {
	AddedItem   Item
	ChangedItem Item
	RemovedItem Item
}

type Item struct {
	Path string
	Size int64
}

type Totals struct {
	AddedSize   int64
	ChangedSize int64
	RemovedSize int64
}

func decode(r io.Reader) (TMUtilCompare, error) {
	buf, err := plist.ToJSON(r)
	if err != nil {
		return TMUtilCompare{}, err
	}
	var result TMUtilCompare
	err = json.Unmarshal(buf, &result)
	return result, err
}
