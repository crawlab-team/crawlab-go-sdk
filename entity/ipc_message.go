package entity

import (
	"encoding/json"
	"github.com/apex/log"
)

type IPCMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	IPC     bool        `json:"ipc"`
}

func (msg *IPCMessage) ToJSON() string {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("failed to marshal IPC message: %v", err)
		return ""
	}
	return string(data)
}
