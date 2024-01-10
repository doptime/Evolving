package main

import (
	"Evolving/gpt"
	"fmt"
	"strings"

	"github.com/yangkequn/saavuu/data"
)

var keyProductLog = data.New[string, string]("ProdLog")

func MMUToTasks(projectName, Description string) {
	var (
		asw *gpt.OutChatGPT
		err error
		jss []string
	)
	ProjectLogger := keyProductLog.Concat(projectName)

	prompt := fmt.Sprintf(`作为Evolute system by maximization of marginal utility代理，你的任务是通过最大化效用的原则来演进一个系统。待演进的产品名称：%s，描述：%s。

	角色一：产品经理
	任务：设计、开发或改进产品。步骤1.通过执行linux 命令或是提出相关问题，建立关于产品的足够必要的认知. 2.根据需求分解业务规格，并评估作为改进方案的业务规格的边际效用。提出并保留具有最高边际效用的一条业务规格
	操作：
	1. 执行Linux命令查看文件，例如：cat filexxx, ls -l, ...  执行linux cmd格式是:
	$_json
	{
		"cmd": "xxx"
	}$_。
	2. 提出项目相关问题，格式为：
		$_json
		{"myQuestion": "xxx"}
		$_
	3. 生成具有最高边际效用的业务规格列表。你需要用迭代的方式做到这一点。首先需要，你出改进一个目标的具有最高边际效用的业务规格，然后提出潜在的更具有边际效用的业务规格，再重复一次这个过程直到确定最终的一条业务规格。格式为：
	 $_json
	 {
	 "BusinessSpecificationItemOfMaxMarginalUtilityVersion1": { "name":  "xxx", "description": "xxx", },
	 "reason-version1-not-maximized-of-marginal-utility": "如果事实上存在比他更具边际效用的业务规格需求，那应该是xxx",
	 "BusinessSpecificationItemOfMaxMarginalUtilityVersion2": { "name":  "xxx", "description": "xxx", },
	 "reason-version2-not-maximized-of-marginal-utility": "如果事实上存在比他更具边际效用的业务规格需求，那应该是xxx",
	 "BusinessSpecificationItemOfMaxMarginalUtilityFinal": { "name":  "xxx", "description": "xxx", },
	 "reason-version1-not-maximized-of-marginal-utility": "如果事实上存在比他更大的效用，那应该是xxx",
	 }
	 $_
	角色二：架构师
	任务：根据产品经理生成的业务规格要求，拆分为技术要求。you need to build the tech requirement in a evolution way, by generate the tech requrement items, and point out the drawback, and then generate the next version of tech requirement items. so on to get the final version of tech requirement items. the format of tech requirement items is:
		$_json
		{
		   "TechRequirementVersion1": [
			   { "name": "xxx", "description": "xxx" },
			   ...
		   ],
		   "DrawbackOfVersion1": "xxx",
		   "TechRequirementVersion2": [
			   { "name": "xxx", "description": "xxx" },
			   ...
		   ],
		   ,"DrawbackOfVersion2": "xxx",
		   "FinalTechRequirement": [
			   { "name": "xxx", "description": "xxx" },
			   ...
		   ],
		
		}$_

	角色三：技术领导
	任务：根据架构师列出的最终技术要求，拆分为开发任务。
	操作：
	1. 生成开发任务版本，指出不足，进化至最终任务列表，格式为：
		$_json
		{
			"DevelopmentTasks": [
				{ "name": "xxx", "description": "xxx" },
				...
			],
			"Drawback": "xxx"
		}
		$_
`, projectName, Description)

	prompt = strings.Replace(prompt, "$_", "```", -1)
	Model := []string{"gpt-3.5-turbo-1106", "gpt-4-plugins", "gpt-4-gizmo"}[1]
	JsonCompleted := gpt.JsonCompleted(1)
	for finish := false; !finish; {
		ProjectLogger.LPush(prompt)
		fmt.Println("product manager is given the prompt:\n", prompt)
		if asw, err = gpt.ApiChatGptXY(prompt, Model, 1, JsonCompleted); err != nil {
			fmt.Println(err)
		}
		if asw == nil || len(asw.Answer) == 0 {
			fmt.Println("no answer")
			break
		}
		ProjectLogger.LPush(asw.Answer)
		jss, _ = JsonCompleted(asw.Answer)
		for _, js := range jss {
			if strings.Contains(js, "myQuestion") {
				if prompt, err = askQuestion(js); err != nil {
					fmt.Println(err)
					return
				}
				break
			} else if strings.Contains(js, "cmd") {
				if prompt, err = RunJSCmd(js); err != nil {
					fmt.Println(err)
					return
				}
				break
			} else if strings.Contains(js, "SpecificationName") {
				Step1ProcessSpecification(projectName, js)
			} else if strings.Contains(js, "FinalTechRequirement") {
				SaveFinalTechRequirement(projectName, js)
			} else if strings.Contains(js, "FinalDevelopmentTasks") {
				SaveFinalDevelopmentTasks(projectName, js)
			}

		}

	}
	// set the standard input, output, and error to the current process's standard input, output, and error

}

// prompt_templete:=`you are a agent named MUEvo. you will evolute a system using maximization of utility. 需要演进的产品名称是: %s, 描述是%s。
// your will fulfill the goal by the following steps:
// 1st. build the one specification with maximization of utility to improve the product.
//    在这里你将作为一个产品经理设计开发或者改进一个产品。现在请你查询必要的信息，以便把这个需求分解成业务规格列表,并从这个规格列表。你可以通过执行linux cmd 命令查看文件内容。比如:
//    通过目标各种 可能的改进方案，并评估改进方案的边际效用。 来设计开发或者改进一个产品。现在请你查询必要的信息，以便把这个需求分解成业务规格列表。你可以通过执行linux cmd 命令查看文件内容。比如:
//   cat filexxx
//   也可以提出关于这个项目的疑问，以便你可以完成创建业务规格列表,业务规格列表的格式是:
//    其中，Description 的描述必须是完备的。也就是说，后续仅仅依靠这个描述，就可以完成技术规格的设计和开发。

//    提问的格式ui是:
//    $_json
//    {
// 	   "myQuestion": "xxx"
//    }$_

//    执行linux cmd格式是:
//    $_json
//    {
// 	   "cmd": "xxx"
//    }$_。
//    最后你需要用迭代的方式，生成一个具有最高边际效用的需求规格说明。
//    首先需要，你出改进一个目标的具有最高边际效用的方案，然后你需要指出reason-not-maximized-of-utility。
//    然后重复下面这个流程2次，以便得到最高边际效用的一个改进方案:把这个潜在的不足考虑在内，变成一个新的，具有更具体，更高边际效用的方案。生成的的格式就像这样:
// $_json
// {
// "BusinessSpecificationItemOfMaxMarginalUtilityVersion1": { "name":  "xxx", "description": "xxx", },
// "reason-version1-not-maximized-of-marginal-utility": "如果事实上存在比他更大的效用，那应该是xxx",
// "BusinessSpecificationItemOfMaxMarginalUtilityVersion2": { "name":  "xxx", "description": "xxx", },
// "reason-version2-not-maximized-of-marginal-utility": "如果事实上存在比他更大的效用，那应该是xxx",
// "BusinessSpecificationItemOfMaxMarginalUtilityFinal": { "name":  "xxx", "description": "xxx", },
// "reason-version1-not-maximized-of-marginal-utility": "如果事实上存在比他更大的效用，那应该是xxx",
// }
// $_
// 2nd. you are a Architect.  according to the specification requirement generated in last step  by product manager , you are going to break down it into tech requirements 。 you need to build the tech requirement in a evolution way, by generate the tech requrement items, and point out the drawback, and then generate the next version of tech requirement items. so on to get the final version of tech requirement items. the format of tech requirement items is:
// $_json
// {
//    "TechRequirementNameVersion1": [
// 	   { "name": "xxx", "description": "xxx" },
// 	   ...
//    ],
//    "DrawbackOfVersion1": "xxx",
//    "TechRequirementNameVersion2": [
// 	   { "name": "xxx", "description": "xxx" },
// 	   ...
//    ],
//    ,"DrawbackOfVersion2": "xxx",
//    "FinalTechRequirement": [
// 	   { "name": "xxx", "description": "xxx" },
// 	   ...
//    ],

// }$_
// 3rd. now you are a Tech leader. According to a FinalTechRequirement listed from Architect agent, you are going to break down into development tasks。业务规格列表的格式是:
// $_json
// {
//    "DevelopmentTasksVersion1": [
// 	   { "name": "xxx", "description": "xxx" },
// 	   ...
//    ],
//    "DrawbackOfVersion1": "xxx",
//    "DevelopmentTasksVersion2": [
// 	   { "name": "xxx", "description": "xxx" },
// 	   ...
//    ],
//    ,"DrawbackOfVersion2": "xxx",
//    "FinalDevelopmentTasks": [
// 	   { "name": "xxx", "description": "xxx" },
// 	   ...
//    ],

// }
// `
