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
