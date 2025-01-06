package crawlab

import (
	"github.com/apex/log"
	"github.com/gocolly/colly/v2"
)

func CollyOnHTMLOne(c *colly.Collector, goqueryString string, getItems func(element *colly.HTMLElement) map[string]any) {
	c.OnHTML(goqueryString, func(element *colly.HTMLElement) {
		item := getItems(element)
		err := SaveItem(item)
		if err != nil {
			log.Errorf("error saving items: %v", err)
		}
	})
}

func CollyOnHTMLMany(c *colly.Collector, goqueryString string, getItems func(element *colly.HTMLElement) []map[string]any) {
	c.OnHTML(goqueryString, func(element *colly.HTMLElement) {
		items := getItems(element)
		err := SaveItem(items...)
		if err != nil {
			log.Errorf("error saving items: %v", err)
		}
	})
}
