package main

import (
	"encoding/json"
	"fmt"

	"github.com/yangkequn/saavuu/data"
)

type SpecificationItem struct {
	Name        string `json:"BusinessSpecificationNameForMaxMarginalUtilityItem"`
	Description string `json:"description"`
}
type MyQuestion struct {
	MyQuestion string `json:"myQuestion"`
}

func KeyMMUEvoSpecifications(projectName string) *data.Ctx[string, *SpecificationItem] {
	return data.New[string, *SpecificationItem]("MMUEvo").Concat(projectName)
}
func Step1ProcessSpecification(projectName, js string) (err error) {

	specificationItem := &SpecificationItem{}
	if err = json.Unmarshal([]byte(js), &specificationItem); err != nil {
		return fmt.Errorf("ProcessSpecification: %w", err)
	}
	KeyMMUEvoSpecifications(projectName).HSet("Specification", specificationItem)
	return
}
