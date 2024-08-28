package handler

import (
	"encoding/json"
	"im/ws"
)

// 解析 Data 到指定的结构体
func ParseData(message *ws.Message, target interface{}) error {
	// 将 Data 字段转换为 JSON 字符串
	dataJSON, err := json.Marshal(message.Data)
	if err != nil {
		return err
	}
	// 解析 JSON 字符串到目标结构体
	err = json.Unmarshal(dataJSON, target)
	if err != nil {
		return err
	}
	return nil
}
