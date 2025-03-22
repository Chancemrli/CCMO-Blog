package dify

import "fmt"

type DifyClient struct {
	addr string
}

func (d *DifyClient) Workflow(data map[string]interface{}) (any, error) {
	fmt.Println("Im dify")
	return nil, nil
}
