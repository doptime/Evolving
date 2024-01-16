package main

import (
	"MMUEvo/apiChatGPT"
	"encoding/json"

	"github.com/yangkequn/saavuu/data"
)

func KeyMMUEvoBusinessState(projectName string) *data.Ctx[string, string] {
	return data.New[string, string]("BusinessState").Concat(projectName)
}

type BusinessStateDescription struct {
	BusinessStateDescription string `json:"BusinessStateDescription"`
}

func UpdateBusinessStateSuccess(projectName string, code *apiChatGPT.GPTResponseCode) (ok bool) {
	bs := &BusinessStateDescription{}
	if code.Type != "json" || json.Unmarshal([]byte(code.Text), bs) != nil || len(bs.BusinessStateDescription) == 0 {
		return false
	} else {
		dataBusinessState := KeyMMUEvoBusinessState(projectName)
		dataBusinessState.Set("BusinessStateDescription", bs.BusinessStateDescription, 0)
		return true
	}
}
