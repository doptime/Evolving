package main

import (
	"encoding/json"
	"fmt"

	"github.com/yangkequn/saavuu/data"
)

type SpecificationItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
type MyQuestion struct {
	MyQuestion string `json:"myQuestion"`
}

func KeyMMUEvoSpecifications(projectName string) *data.Ctx[string, *SpecificationItem] {
	return data.New[string, *SpecificationItem]("MMUEvo").Concat(projectName)
}
func Step1ProcessSpecification(projectName, js string) (err error) {
	type SpecificationVersion struct {
		Version           int    `json:"version"`
		SpecificationName string `json:"specificationName"`
		Description       string `json:"description"`
		Reason            string `json:"reason-not-maximized-of-utility"`
	}

	SpecificationVersions := []*SpecificationVersion{}
	if err = json.Unmarshal([]byte(js), &SpecificationVersions); err != nil {
		fmt.Println(err)
		return
	} else if len(SpecificationVersions) == 0 {
		fmt.Println("no SpecificationVersions")
		return
	}
	SpecificationLastVersion := SpecificationVersions[len(SpecificationVersions)-1]
	SpecificationItem := &SpecificationItem{
		Name:        SpecificationLastVersion.SpecificationName,
		Description: SpecificationLastVersion.Description,
	}
	KeyMMUEvoSpecifications(projectName).HSet("Specification", SpecificationItem)
	return
}

// var keySpecificationList = data.New[string, []*SpecificationItem]("SpecificationList")
// var keyProductLog = data.New[string, string]("ProdLog")

// func ProductManager_BuildSpecificationList(SkillTreeName string) {
// 	var (
// 		asw               *gpt.OutChatGPT
// 		err               error
// 		specificationList *SpecificationList
// 		myQuestion        *MyQuestion
// 		cmd               *Command
// 		jss               []string
// 	)
// 	// the first argument is the command
// 	projectName := "youtube节目改进,主题：对话节目"
// 	PMLogger := keyProductLog.Concat("Manager").Concat(projectName)
// 	Description := "你将要改一个Youtube节目。是个对话节目。你需要改进这个节目使得内容更具吸引力。这个节目包含一段对话, 位于./dialouge.txt。还有针对这个对话的评论,位于/commend.txt。你同样需要改进这个评论使得更具有吸引力。"
// 	prompt := fmt.Sprintf(` you are a product manager agent powered by LLM. 你的作用是通过考虑一个目标可能的改进方案，评估改进方案的边际效用。 来设计开发或者改进一个产品。产品名称是%s。产品描述是%s。 现在请你查询必要的信息，以便把这个需求分解成业务规格列表。你可以通过执行linux cmd 命令查看文件内容。比如:
// 	cat filexxx
// 	也可以提出关于这个项目的疑问，以便你可以完成创建业务规格列表,业务规格列表的格式是:
// 	$_json
// 	{
// 		"SpecificationItems": [
// 			{ "name": "xxx", "Description": "xxx" },
// 			...
// 		]
// 	}$_
// 	其中，Description 的描述必须是完备的。也就是说，后续仅仅依靠这个描述，就可以完成技术规格的设计和开发。

// 	提问的格式ui是:
// 	$_json
// 	{
// 		"myQuestion": "xxx"
// 	}$_

// 	执行linux cmd格式是:
// 	$_json
// 	{
// 		"cmd": "xxx"
// 	}$_
// 	`, projectName, Description)

// 	prompt = strings.Replace(prompt, "$_", "```", -1)
// 	Model := []string{"gpt-3.5-turbo-1106", "gpt-4-plugins", "gpt-4-gizmo"}[1]
// 	JsonCompleted := gpt.JsonCompleted(1)
// 	for finish := false; !finish; {
// 		PMLogger.LPush(prompt)
// 		fmt.Println("product manager is given the prompt:\n", prompt)
// 		if asw, err = gpt.ApiChatGptXY(prompt, Model, 1, JsonCompleted); err != nil {
// 			fmt.Println(err)
// 		}
// 		if asw == nil || len(asw.Answer) == 0 {
// 			fmt.Println("no answer")
// 			break
// 		}
// 		PMLogger.LPush(asw.Answer)
// 		jss, _ = JsonCompleted(asw.Answer)
// 		for _, js := range jss {
// 			if err = json.Unmarshal([]byte(js), &specificationList); err == nil {
// 				keySpecificationList.Set(projectName, specificationList.SpecificationItems, 0)
// 				finish = true
// 				break
// 			}
// 			if err = json.Unmarshal([]byte(js), &myQuestion); err == nil {
// 				//print the question to the console, and wait for the answer
// 				fmt.Println(" the product manager has a question: " + myQuestion.MyQuestion)
// 				//read the answer from the console
// 				fmt.Scanln(&prompt)
// 				break
// 			}
// 			if err = json.Unmarshal([]byte(js), &cmd); err == nil {
// 				//print the question to the console, and wait for the answer
// 				fmt.Println(" the product manager has a question: " + cmd.Command)
// 				//read the answer from the console
// 				commandStrings := strings.Split(cmd.Command, " ")
// 				if len(commandStrings) == 0 {
// 					finish = true
// 					fmt.Println("no command")
// 					break
// 				}
// 				prompt = RunCmd(commandStrings[0], commandStrings[1:]...)
// 				break
// 			}
// 		}

// 	}
// 	// set the standard input, output, and error to the current process's standard input, output, and error

// }
