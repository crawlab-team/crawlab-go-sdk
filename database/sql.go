package database

import (
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-go-sdk/constants"
	"github.com/crawlab-team/crawlab-go-sdk/entity"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"os"
	"runtime/debug"
	"strings"
)

var Db *sqlx.DB

func GetSqlDatabaseConnectionString(dataSourceType string, host string, port string, username string, password string, database string) (connStr string, err error) {
	// 获取数据库链接串
	if dataSourceType == constants.DataSourceTypeMysql {
		connStr = fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, database)
	} else if dataSourceType == constants.DataSourceTypePostgres {
		connStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", host, port, username, database, password, "disable")
	} else {
		err = errors.New(dataSourceType + " is not implemented")
		log.Errorf(err.Error())
		debug.PrintStack()
		return connStr, err
	}
	return connStr, nil
}

func GetSqlDatabase(ds entity.DataSource) *sqlx.DB {
	if Db != nil {
		return Db
	}

	// 获取数据库链接串
	connStr, err := GetSqlDatabaseConnectionString(ds.Type, ds.Host, ds.Port, ds.Username, ds.Password, ds.Database)
	if err != nil {
		log.Errorf("get connection string error: " + err.Error())
		debug.PrintStack()
		return Db
	}

	// 数据库
	Db, err = sqlx.Open(ds.Type, connStr)
	if err != nil {
		log.Errorf("open database error: " + err.Error())
		debug.PrintStack()
	}

	return Db
}

func InsertItem(ds entity.DataSource, item entity.Item) (err error) {
	db := GetSqlDatabase(ds)
	col := os.Getenv("CRAWLAB_COLLECTION")
	if _, err := db.Exec(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		col,
		strings.Join(GetItemKeys(item), ","),
		strings.Join(GetItemValuesWithQuote(item), ","),
	)); err != nil {
		log.Errorf("insert item error: " + err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

func UpdateItem(ds entity.DataSource, item entity.Item, dedupField string) (err error) {
	db := GetSqlDatabase(ds)
	col := os.Getenv("CRAWLAB_COLLECTION")
	if _, err := db.Exec(fmt.Sprintf("UPDATE %s SET %s WHERE %s = '%s'",
		col,
		GetItemSqlUpdateStr(item),
		dedupField,
		item[dedupField],
	)); err != nil {
		log.Errorf("update item error: " + err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

func GetItem(ds entity.DataSource, key string, value string) (item entity.Item, err error) {
	db := GetSqlDatabase(ds)
	col := os.Getenv("CRAWLAB_COLLECTION")
	rows, err := db.Queryx(fmt.Sprintf("SELECT * FROM %s WHERE %s = '%s' LIMIT 1",
		col,
		key,
		value,
	))
	if err != nil {
		log.Errorf("get item error: " + err.Error())
		debug.PrintStack()
		return item, err
	}
	if rows.Next() {
		item = make(map[string]interface{})
		if err := rows.MapScan(item); err != nil {
			log.Errorf("get item error: " + err.Error())
			debug.PrintStack()
			return item, err
		}

		// 字节类型转化
		// Mysql 默认是字节类型, 需要转化为字符串
		for k, v := range item {
			switch v.(type) {
			case []byte:
				item[k] = string(v.([]byte))
			}
		}
	}
	return item, nil
}

func DeleteItems(ds entity.DataSource, key string, value string) (err error) {
	db := GetSqlDatabase(ds)
	col := os.Getenv("CRAWLAB_COLLECTION")
	_, err = db.Exec(fmt.Sprintf("DELETE FROM %s WHERE %s = '%s'",
		col,
		key,
		value,
	))
	if err != nil {
		log.Errorf("delete items error: " + err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}

func GetItemKeys(item entity.Item) (res []string) {
	for k := range item {
		res = append(res, k)
	}
	return res
}

func GetItemValues(item entity.Item) (res []string) {
	for _, v := range item {
		res = append(res, v.(string))
	}
	return res
}

func GetItemValuesWithQuote(item entity.Item) (res []string) {
	for _, v := range item {
		res = append(res, "'"+v.(string)+"'")
	}
	return res
}

func GetItemSqlUpdateStr(item entity.Item) string {
	var arr []string
	for k, v := range item {
		arr = append(arr, fmt.Sprintf("%s = '%s'", k, v))
	}
	return strings.Join(arr, ",")
}
