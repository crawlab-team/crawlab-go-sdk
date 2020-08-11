package database

import (
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/crawlab-team/crawlab-go-sdk/entity"
	"github.com/segmentio/kafka-go"
	"runtime/debug"
)

func GetKafkaConnection(ds entity.DataSource) (conn *kafka.Conn, err error) {
	conn, err = kafka.Dial(
		"tcp",
		fmt.Sprintf("%s:%s", ds.Host, ds.Port),
	)
	if err != nil {
		log.Errorf("dial kafka error: " + err.Error())
		debug.PrintStack()
		return conn, err
	}
	return conn, nil
}

func SendKafkaMsg(ds entity.DataSource, item entity.Item) (err error) {
	conn, err := GetKafkaConnection(ds)
	if err != nil {
		return err
	}
	msgStr, err := json.Marshal(&item)
	if err != nil {
		log.Errorf("marshal json error: " + err.Error())
		debug.PrintStack()
		return err
	}
	if _, err := conn.WriteMessages(
		kafka.Message{
			Topic:     ds.Database,
			Partition: 0,
			Value:     msgStr,
		},
	); err != nil {
		log.Errorf("send kafka message error: " + err.Error())
		debug.PrintStack()
		return err
	}
	return nil
}
