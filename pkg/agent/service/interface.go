package service

import "fmt"

type AgentService interface {
	Workflow(data map[string]interface{}) (any, error)
}

func NewAgent(agent AgentService) {
	fmt.Println(agent.Workflow(map[string]interface{}{
		"input": "Hello, World!",
	}))
}
