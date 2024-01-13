package main

import (
	"encoding/json"
	"strings"

	"github.com/yangkequn/saavuu/data"
)

func KeyMMUEvoBusinessState(projectName string) *data.Ctx[string, string] {
	return data.New[string, string]("BusinessState").Concat(projectName)
}

func UpdateBusinessState(projectName, js string) (ok bool, err error) {
	type BusinessStateDescription struct {
		BusinessStateDescription string `json:"BusinessStateDescription"`
	}
	if !strings.Contains(js, "BusinessState") {
		return
	}

	bs := &BusinessStateDescription{}
	if err = json.Unmarshal([]byte(js), &bs); err != nil {
		return
	}
	dataBusinessState := KeyMMUEvoBusinessState(projectName)
	err = dataBusinessState.Set("BusinessStateDescription", bs.BusinessStateDescription, 0)
	return true, nil
}
