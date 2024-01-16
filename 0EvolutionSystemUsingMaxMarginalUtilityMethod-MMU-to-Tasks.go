package main

import (
	"MMUEvo/apiChatGPT"
	"MMUEvo/runSysCmd"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/yangkequn/saavuu/data"
)

var keyProductLog = data.New[string, []string]("ProdLog")

func MMUToTasks(projectName, Description string) {
	var (
		asw         *apiChatGPT.OutChatGPT
		err         error
		gptCodes    []*apiChatGPT.GPTResponseCode
		gpt4Session = apiChatGPT.NewChatSession(apiChatGPT.Models["gpt-4-plugins"], 1)
	)
	ProjectLogger := keyProductLog.Concat(projectName)
	BusinessState, _ := KeyMMUEvoBusinessState(projectName).HGet("BusinessStateDescription")

	prompt := fmt.Sprintf(`you are agent to evolute a system by method of maximization of marginal utility，你的任务是通过最大化效用的原则来演进一个系统。注意，你一次只需要关注一小部分的演进目标。当前的演进过程将会被迭代成百上千次。所以你不应该尝试完成大的目标。你需要高质量完成很少的一点目标。
待演进的产品/系统名称：%s
待演进的产品/系统描述：%s
这是上一次迭代生成的 BusinessStateDescription :
%s

你需要在接下来的5个步骤中完成这个任务：


步骤1, 可选.如果必要，你可以通过以下两个步骤进一步了解当前系统的业务状态。
	- 你可以通过执行linux 命令，建立关于产品当前状态的足够必要的认知.
	首先，优先通过使用Linux命令查看文件来了解现有系统的实现。例如：
		- 查看目录下的文件：ls -lh | awk '{print $2,$5,$6,$7,$8,$9}'
		- 查看文件内容以判断是否需要改进(使用 cat -n xxx 命令，因为后续可能要依照行号完成增量修改)：cat -n xxx
	执行linux cmd格式是:
$_bash
# Query Bash Cmds
...
$_。
 	-你可以进一步提出关于该系统的关系核心问题，相应的回答会出现在在下批次的对话的待演进的产品/系统描述中。格式为：
$_json
{"Questions-that-require-further-clarification": "xxx"}
$_

步骤2：你现在的角色是设计、开发或改进产品的产品经理。你需要根据需求分解业务规格，并评估作为改进方案的业务规格的边际效用。提出并生成具有最高边际效用的一条业务规格。你需要用迭代的方式做到这一点。
	- 你提出一个用来改进目标的具有最高边际效用的业务规格
	- 然后讨论改进这个业务规格的边际效用的潜在办法，
	- 并生成最终的业务规格。格式为：
$_json
{
"BusinessSpecificationNameForMaxMarginalUtilityItem":  "xxx", "description": "xxx", 
}$_


步骤3：你现在的角色是技术领导(Tech Lead)。根据刚刚产品经理生成的业务规格，你需要从中确定一条最重要的，未完成的开发任务，作为本次开发目标。注意，这个目标必须是个很小的目标，或者是目标中的一小部分。其最终实现应该能在2K左右的字符内完成。
这是技术要求定案格式:
$_json
{ "DevelopmentTasks": "xxx", "description": "xxx" }
$_

步骤4：你现在的角色是Developer。根据技术领导生成的开发任务，你需要生成对步骤1中查看的文件的的修改。这个修改通过bash命令进行。对增量修改情形，建议采用sed命令，也就是删除行和添加行来实现。所有的linux命令都是可用的，比如touch新建文件。格式为：
$_bash
# Developer Bash Cmds
...
$_

步骤5, 可选：完成必要的对业务系统的理解后，生成内容更新后业务状态的描述。格式为：
$_json
{	"BusinessStateDescription": "xxx"
}$_

现在请按步骤开始你演进系统的工作。你不能在一个回复中生成全部的结果，因为你需要获取足够的信息来理解系统，做出优秀的决策。如果你需要进一步了解现有的实现，请按步骤1进行。

这是待演进的产品/系下运行linux命令 ls -l | awk '{print $2,$5,$6,$7,$8,$9}'的返回结果：
%s
`, projectName, Description, BusinessState, runSysCmd.RunCmd(`ls -lh | awk '{print $2,$5,$6,$7,$8,$9}'`))

	for prompt = strings.Replace(prompt, "$_", "```", -1); prompt != ""; {
		fmt.Println(prompt)
		if asw, err = gpt4Session.ChatOnce(prompt); err != nil {
			fmt.Println(err)
			return
		} else if asw == nil || len(asw.Answer) == 0 {
			log.Fatal().Msgf("no answer: %s", prompt)
		} else {
			ProjectLogger.LPush([]string{prompt, asw.Answer})
		}
		if gptCodes = apiChatGPT.GptResponseParseToCodes(asw.Answer); len(gptCodes) == 0 {
			log.Fatal().Msgf("no answer: %s", prompt)
		}
		prompt = ""
		for _, _code := range gptCodes {
			if UpdateBusinessStateSuccess(projectName, _code) {
			} else if AskHumanHelpInNextIter(_code, ProjectLogger) {
			} else if _code.Type == "bash" && (strings.Contains(_code.Text, "Query") || strings.Contains(_code.Text, "查询")) {
				prompt = runSysCmd.RunCmd(_code.Text)
			} else if RedisSave_BussinessSpecification(projectName, _code) {
			} else if SaveDevelopmentTask(projectName, _code) {
			} else if _code.Type == "bash" && strings.Contains(_code.Text, "Developer Bash Cmds") {
				runSysCmd.RunCmd(_code.Text)
			} else {
				log.Info().Msgf("not  handled: %s", _code.Text)
			}
		}
	}

}
