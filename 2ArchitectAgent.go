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
func SaveFinalTechRequirement(projectName, js string) (err error) {
	var (
		finalTechRequirement     = []*TechRequirementItem{}
		FinalTechRequirementJson string
	)
	gjs := gjson.Parse(js)
	FTR := gjs.Get("FinalTechRequirement")
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

// var keyTechRequirement = data.New[string, []*TechRequirementItem]("TechRequirement")

// func Architect_Breakdown(projectName string) {
// 	var (
// 		asw                 *gpt.OutChatGPT
// 		err                 error
// 		TechRequirementList *TechRequirementList
// 		jss                 []string
// 		specificationItems  []byte
// 		prompt              string
// 		SpecificationItems  []*SpecificationItem
// 	)
// 	if SpecificationItems, err = keySpecificationList.Get(projectName); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	specificationItems, _ = json.Marshal(SpecificationItems)
// 	// the first argument is the command
// 	ArchitectLogger := keyProductLog.Concat("Architect").Concat(projectName)
// 	prompt = fmt.Sprintf(` you are a Architect. given a specification list from product manager, you are going to break down tech requirements according to it 。 产品名称是%s。规格列表是%s。 现在请你查询必要的信息，以便把这个需求分解成业务规格列表。
// 	业务规格列表的格式是:
// 	$_json
// 	{
// 		"TechRequirementItem": [
// 			{ "name": "xxx", "Description": "xxx" },
// 			...
// 		]
// 	}$_
// 	`, projectName, string(specificationItems))

// 	prompt = strings.Replace(prompt, "$_", "```", -1)
// 	Model := gpt.Models["4p"]
// 	JsonCompleted := gpt.JsonCompleted(1)
// 	for finish := false; !finish; {
// 		ArchitectLogger.LPush(prompt)
// 		fmt.Println("product manager is given the prompt:\n", prompt)
// 		if asw, err = gpt.ApiChatGptXY(prompt, Model, 1, JsonCompleted); err != nil {
// 			fmt.Println(err)
// 		}
// 		if asw == nil || len(asw.Answer) == 0 {
// 			fmt.Println("no answer")
// 			break
// 		}
// 		ArchitectLogger.LPush(asw.Answer)
// 		jss, _ = JsonCompleted(asw.Answer)
// 		for _, js := range jss {
// 			if err = json.Unmarshal([]byte(js), &TechRequirementList); err == nil {
// 				keyTechRequirement.Set(projectName, TechRequirementList.TechRequirementItem, 0)
// 				finish = true
// 				break
// 			}
// 		}

// 	}
// 	// set the standard input, output, and error to the current process's standard input, output, and error

// }
