package main

import (
	"CCMO/gozero/blog/pkg/agent/coze"
	"fmt"
)

const (
	PAT     = "pat_47yGCrMq89u18fEwLXsVawDPArWPBHgVyzEnJt2YSXwearg68LQ2T57Zl04EdqSW"
	WID     = "7484100103276740660"
	CONTENT = "好想摆烂了啊，累了"
)

func main() {
	// // 获取连接
	// c := coze.NewCozeClient(PAT)

	// // 工作流Input
	// data := map[string]interface{}{
	// 	"input": CONTENT,
	// }

	// // 工作流请求
	// req := &coze.RunWorkflowsReq{
	// 	WorkflowID: WID,
	// 	Parameters: data,
	// 	// 是否异步执行
	// 	IsAsync: false,
	// }

	// // 建立请求
	// resp, err := c.Client.Workflows.Runs.Create(context.Background(), req)
	// if err != nil {
	// 	fmt.Println("Error running workflow:", err)
	// 	return
	// }
	// fmt.Println(resp)
	// fmt.Println("ok")

	client := coze.NewCozeClient(PAT)
	data, err := client.Workflow(map[string]interface{}{"input": CONTENT}, WID)
	if err != nil {
		fmt.Println("Error running workflow:", err)
		return
	}
	fmt.Println(data)
}
