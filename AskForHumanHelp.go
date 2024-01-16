package main

import (
	"MMUEvo/apiChatGPT"
	"encoding/json"

	"github.com/yangkequn/saavuu/data"
)

type MyQuestion struct {
	MyQuestion string `json:"Questions-that-require-further-clarification"`
}

func AskHumanHelpInNextIter(code *apiChatGPT.GPTResponseCode, ProjectLogger *data.Ctx[string, []string]) (ok bool) {
	var (
		myQuestion *MyQuestion
	)
	if code.Type != "json" || json.Unmarshal([]byte(code.Text), &myQuestion) != nil || len(myQuestion.MyQuestion) == 0 {
		return false
	}
	ProjectLogger.LPush([]string{myQuestion.MyQuestion})
	return true
}
