package utils

import (
	"github.com/crawlab-team/crawlab-go-sdk/database"
	"github.com/crawlab-team/crawlab-go-sdk/entity"
	"github.com/globalsign/mgo/bson"
	"github.com/smartystreets/goconvey/convey"
	"os"
	"testing"
)

func TestSaveItem(t *testing.T) {
	convey.Convey("Test SaveItem", t, func() {
		taskId := "test_task_id"
		url := "http://example.com"
		title := "test"
		_ = os.Setenv("CRAWLAB_TASK_ID", taskId)
		_ = os.Setenv("CRAWLAB_MONGO_HOST", "localhost")
		_ = os.Setenv("CRAWLAB_MONGO_PORT", "27017")
		_ = os.Setenv("CRAWLAB_MONGO_DATABASE", "crawlab_test")
		_ = os.Setenv("CRAWLAB_COLLECTION", "result_test")
		item := entity.Item{
			"url":   url,
			"title": title,
		}
		if err := SaveItem(item); err != nil {
			t.Fatal("save item failed")
		}
		_, c, _ := database.GetMongoCol(entity.DataSource{
			Type: "",
		})
		//defer s.Close()
		var itemDb entity.Item
		_ = c.Find(bson.M{"url": url}).One(&itemDb)
		convey.Convey("title should be 'test'", func() {
			convey.So(itemDb["title"], convey.ShouldEqual, title)
		})
		count, _ := c.Find(bson.M{"url": url}).Count()
		convey.Convey("count should be greater than 0", func() {
			convey.So(count, convey.ShouldBeGreaterThan, 0)
		})
		_, _ = c.RemoveAll(bson.M{"url": url})
	})
}
