package main

import (
	"github.com/crawlab-team/crawlab-go-sdk/entity"
	"github.com/crawlab-team/crawlab-go-sdk/utils"
)

func SaveItem(item entity.Item) (err error) {
	if err := utils.SaveItem(item); err != nil {
		return err
	}
	return nil
}
