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

func (msg *IPCMessage) ToJSON() (res string, err error) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("failed to marshal IPC message: %v", err)
		return "", err
	}
	return string(data), nil
}
