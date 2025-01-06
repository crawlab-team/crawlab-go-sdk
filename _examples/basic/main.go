package main

import (
	"github.com/crawlab-team/crawlab-go-sdk"
)

func main() {
	item := make(map[string]interface{})
	item["url"] = "http://example.com"
	item["title"] = "hello world"
	_ = crawlab.SaveItem(item)
}
