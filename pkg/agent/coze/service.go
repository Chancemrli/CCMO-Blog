package coze

import (
	"context"
	"fmt"

	"github.com/coze-dev/coze-go"
)

// 同步工作流
func (c *CozeClient) Workflow(data map[string]interface{}, wId string) (any, error) {
	// 建立请求
	req := &coze.RunWorkflowsReq{
		WorkflowID: wId,
		Parameters: data,
		// 是否异步执行
		IsAsync: false,
	}

	// 发送请求
	resp, err := c.Client.Workflows.Runs.Create(context.Background(), req)
	if err != nil {
		fmt.Println("Error running workflow:", err)
		return nil, err
	}

	return resp, nil
}
