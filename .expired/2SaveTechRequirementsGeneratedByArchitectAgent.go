package main

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/yangkequn/saavuu/data"
)

type TechRequirementItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type TechRequirementList struct {
	TechRequirementItem []*TechRequirementItem `json:"SpecificationList"`
}

func KeyMMUEvoTechRequirements(projectName string) *data.Ctx[string, []*TechRequirementItem] {
	return data.New[string, []*TechRequirementItem]("MMUEvo").Concat(projectName)
}
func SaveTechRequirement(projectName, js string) (err error) {
	var (
		finalTechRequirement     = []*TechRequirementItem{}
		FinalTechRequirementJson string
	)
	gjs := gjson.Parse(js)
	FTR := gjs.Get("TechRequirements")
	if !FTR.Exists() || !FTR.IsArray() {
		return fmt.Errorf("no FinalTechRequirement")
	}
	FinalTechRequirementJson = FTR.String()
	if err = json.Unmarshal([]byte(FinalTechRequirementJson), &finalTechRequirement); err != nil {
		return err
	}
	KeyMMUEvoTechRequirements(projectName).HSet("TechRequirement", finalTechRequirement)
	return
}
