package crawlab

import (
	"fmt"
	"github.com/crawlab-team/crawlab-go-sdk/constants"
	"github.com/crawlab-team/crawlab-go-sdk/entity"
)

func SaveItem(items ...map[string]any) {
	SaveItems(items)
}

func SaveItems(items []map[string]any) {
	msg := entity.IPCMessage{
		Type:    constants.IPCMessageTypeData,
		Payload: items,
		IPC:     true,
	}
	jsonData := msg.ToJSON()
	fmt.Println(jsonData)
}
