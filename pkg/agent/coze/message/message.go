package message

import (
	"encoding/json"
	"log"

	"github.com/coze-dev/coze-go"
)

type Data struct {
	ContentType    int    `json:"content_type"`
	Data           string `json:"data"`
	OriginalResult *any   `json:"original_result"`
	TypeForModel   int    `json:"type_for_model"`
}

// 解析主要内容
func ParseData(resp coze.RunWorkflowsResp) (string, error) {
	var dataStruct Data
	err := json.Unmarshal([]byte(resp.Data), &dataStruct)
	if err != nil {
		log.Fatalf("Error unmarshalling Data: %v", err)
		return "", err
	}
	return dataStruct.Data, nil
}
