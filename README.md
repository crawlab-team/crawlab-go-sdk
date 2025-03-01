# Crawlab Go SDK

Crawlab Go SDK supports Golang-based spiders integration with Crawlab. It contains a number of APIs including saving crawled items into different data sources including MongoDB, MySQL, Postgres, ElasticSearch and Kafka.

## Basic Usage

```go
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
```

## Example Using Colly

```go
package main

import (
	"fmt"
	"github.com/crawlab-team/crawlab-go-sdk"

	"github.com/gocolly/colly/v2"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: quotes.toscrape.com
		colly.AllowedDomains("quotes.toscrape.com"),
	)

	// On every a element which has href attribute call callback
	crawlab.CollyOnHTMLMany(c, "a[href]", func(e *colly.HTMLElement) []map[string]any {
		return []map[string]any{
			{
				"text": e.Text,
				"link": e.Attr("href"),
			},
		}
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://quotes.toscrape.com
	c.Visit("https://quotes.toscrape.com")
}
```