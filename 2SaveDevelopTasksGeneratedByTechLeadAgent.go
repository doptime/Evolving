package main

import (
	"MMUEvo/apiChatGPT"
	"encoding/json"

	"github.com/yangkequn/saavuu/data"
)

type DevelopmentTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func KeyMMUEvoDevelopmentTasks(projectName string) *data.Ctx[string, *DevelopmentTask] {
	return data.New[string, *DevelopmentTask]("MMUEvo").Concat(projectName)
}

func SaveDevelopmentTask(projectName string, code *apiChatGPT.GPTResponseCode) (ok bool) {
	var (
		devTasks = []*DevelopmentTask{}
	)

	if code.Type != "json" || json.Unmarshal([]byte(code.Text), &devTasks) != nil {
		return false
	}
	KeyMMUEvoDevelopmentTasks(projectName).HSet("DevelopmentTask", devTasks)
	return
}
