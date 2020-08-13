package utils

import (
	"encoding/json"
	"github.com/crawlab-team/crawlab-go-sdk/constants"
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

func TestSaveItemSqlMySQL(t *testing.T) {
	convey.Convey("Test SaveItemSql (MySQL)", t, func() {
		taskId := "test_task_id"
		url := "http://example.com"
		title := "test"
		ds := entity.DataSource{
			Type:     constants.DataSourceTypeMysql,
			Host:     "localhost",
			Port:     "3306",
			Database: "test",
			Username: "root",
			Password: "mysql",
		}
		dsBytes, _ := json.Marshal(&ds)
		dsStr := string(dsBytes)
		_ = os.Setenv("CRAWLAB_TASK_ID", taskId)
		_ = os.Setenv("CRAWLAB_COLLECTION", "results2")
		_ = os.Setenv("CRAWLAB_DATA_SOURCE", dsStr)
		item := entity.Item{
			"url":   url,
			"title": title,
		}
		if err := SaveItem(item); err != nil {
			t.Fatal("save item failed")
		}
		itemDb, _ := database.GetItem(ds, "url", url)
		convey.Convey("title should be 'test'", func() {
			convey.So(itemDb["title"], convey.ShouldEqual, title)
		})
		_ = database.DeleteItems(ds, "url", url)
	})
}

func TestSaveItemSqlPostgres(t *testing.T) {
	convey.Convey("Test SaveItemSql (Postgres)", t, func() {
		taskId := "test_task_id"
		url := "http://example.com"
		title := "test"
		ds := entity.DataSource{
			Type:     constants.DataSourceTypePostgres,
			Host:     "localhost",
			Port:     "5432",
			Database: "postgres",
			Username: "postgres",
			Password: "postgres",
		}
		dsBytes, _ := json.Marshal(&ds)
		dsStr := string(dsBytes)
		_ = os.Setenv("CRAWLAB_TASK_ID", taskId)
		_ = os.Setenv("CRAWLAB_COLLECTION", "results2")
		_ = os.Setenv("CRAWLAB_DATA_SOURCE", dsStr)
		item := entity.Item{
			"url":   url,
			"title": title,
		}
		if err := SaveItem(item); err != nil {
			t.Fatal("save item failed")
		}
		itemDb, _ := database.GetItem(ds, "url", url)
		convey.Convey("title should be 'test'", func() {
			convey.So(itemDb["title"], convey.ShouldEqual, title)
		})
		_ = database.DeleteItems(ds, "url", url)
	})
}
