package main

import (
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/yangkequn/saavuu/data"
)

type DevelopmentTask struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func KeyMMUEvoDevelopmentTasks(projectName string) *data.Ctx[string, []*DevelopmentTask] {
	return data.New[string, []*DevelopmentTask]("MMUEvo").Concat(projectName)
}

func SaveFinalDevelopmentTasks(projectName, js string) (err error) {
	var (
		devTasks                 = []*DevelopmentTask{}
		FinalTechRequirementJson string
	)
	gjs := gjson.Parse(js)
	FTR := gjs.Get("FinalDevelopmentTasks")
	if !FTR.Exists() || !FTR.IsArray() {
		return fmt.Errorf("no FinalTechRequirement")
	}
	FinalTechRequirementJson = FTR.String()
	if err = json.Unmarshal([]byte(FinalTechRequirementJson), &devTasks); err != nil {
		return err
	}
	KeyMMUEvoDevelopmentTasks(projectName).HSet("DevelopmentTasks", devTasks)
	return
}

// var keyDevelopTasks = data.New[string, []*DevelopmentTask]("DevelopTasks")

// func TechLead_DevelopmentBreakdown(projectName string) {
// 	var (
// 		asw                  *gpt.OutChatGPT
// 		err                  error
// 		TechRequirementItems []*TechRequirementItem
// 		jss                  []string
// 		techRequirementItems []byte
// 		developTaskList      *DevelopTaskList
// 		prompt               string
// 	)
// 	if TechRequirementItems, err = keyTechRequirement.Get(projectName); err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	techRequirementItems, _ = json.Marshal(TechRequirementItems)
// 	// the first argument is the command
// 	ArchitectLogger := keyProductLog.Concat("Architect").Concat(projectName)
// 	prompt = fmt.Sprintf(` you are a Tech leader. According to a tech requirement list from Architect agent, you are going to break down into development tasks。 产品名称是%s。规格列表是%s。 现在请你查询必要的信息，以便把这个需求分解成业务规格列表。
// 	业务规格列表的格式是:
// 	$_json
// 	{
// 		"TechRequirementItem": [
// 			{ "name": "xxx", "Description": "xxx" },
// 			...
// 		]
// 	}$_
// 	`, projectName, string(techRequirementItems))

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
// 			if err = json.Unmarshal([]byte(js), &developTaskList); err == nil {
// 				keyDevelopTasks.Set(projectName, developTaskList.DevelopTasks, 0)
// 				finish = true
// 				break
// 			}
// 		}

// 	}
// 	// set the standard input, output, and error to the current process's standard input, output, and error

// }
