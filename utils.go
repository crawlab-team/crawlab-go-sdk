package crawlab

import (
	"fmt"
	"github.com/crawlab-team/crawlab-go-sdk/constants"
	"github.com/crawlab-team/crawlab-go-sdk/entity"
)

func SaveItem(items ...map[string]any) (err error) {
	return SaveItems(items)
}

func SaveItems(items []map[string]any) (err error) {
	msg := entity.IPCMessage{
		Type:    constants.IPCMessageTypeData,
		Payload: items,
		IPC:     true,
	}
	jsonData, err := msg.ToJSON()
	if err != nil {
		return err
	}
	fmt.Println(jsonData)
	return nil
}
