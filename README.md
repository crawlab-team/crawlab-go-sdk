# Crawlab Go SDK

Crawlab Go SDK supports Golang-based spiders integration with Crawlab. It contains a number of APIs including saving crawled items into different data sources including MongoDB, MySQL, Postgres, ElasticSearch and Kafka.

## Basic Usage

```go
package main

import (
	crawlab "github.com/crawlab-team/crawlab-go-sdk"
)

func main() {
	item := make(map[string]interface{})
	item["url"] = "http://example.com"
	item["title"] = "hello world"
	_ = crawlab.SaveItem(item)
}

```

## Example Using Colly

```go
package main

import (
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-go-sdk"
	"github.com/gocolly/colly/v2"
	"runtime/debug"
)

func main() {
	startUrl := "https://www.baidu.com/s?wd=crawlab"

	c := colly.NewCollector(
		colly.AllowedDomains("www.baidu.com"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.105 Safari/537.36"),
	)

	c.OnHTML("#content_left > .c-container", func(e *colly.HTMLElement) {
		item := make(map[string]interface{})
		item["title"] = e.ChildText("h3.t > a")
		item["url"] = e.ChildAttr("h3.t > a", "href")
		if err := crawlab.SaveItem(item); err != nil {
			log.Errorf("save item error: %v", err)
			debug.PrintStack()
			return
		}
	})

	c.OnRequest(func(r *colly.Request) {
		log.Debugf(fmt.Sprintf("Visiting %s", r.URL.String()))
	})

	if err := c.Visit(startUrl); err != nil {
		log.Errorf("visit error: " + err.Error())
		debug.PrintStack()
		panic(fmt.Sprintf("Unable to visit %s", startUrl))
	}

	c.Wait()
}
```