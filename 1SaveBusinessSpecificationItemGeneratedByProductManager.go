package main

import (
	"MMUEvo/apiChatGPT"
	"encoding/json"

	"github.com/yangkequn/saavuu/data"
)

type SpecificationItem struct {
	Name        string `json:"BusinessSpecificationNameForMaxMarginalUtilityItem"`
	Description string `json:"description"`
}

func KeyMMUEvoSpecifications(projectName string) *data.Ctx[string, *SpecificationItem] {
	return data.New[string, *SpecificationItem]("MMUEvo").Concat(projectName)
}
func RedisSave_BussinessSpecification(projectName string, code *apiChatGPT.GPTResponseCode) (done bool) {
	specificationItem := &SpecificationItem{}
	if code.Type != "json" || json.Unmarshal([]byte(code.Text), &specificationItem) != nil {
		return false
	}
	KeyMMUEvoSpecifications(projectName).HSet("Specification", specificationItem)
	return
}
