package gpt

import "strings"

func JsonCompleted(N int) func(Answer string) (JSParts []string, done bool) {
	extractOneJson := func(Answer *string, AnswerParts *[]string) (success bool) {
		var (
			ind int
		)
		if ind = strings.Index(*Answer, "```json"); ind < 0 {
			return false
		}
		*Answer = (*Answer)[ind+7:]
		if ind = strings.Index((*Answer), "```"); ind < 0 {
			return false
		}
		jsonStr := strings.TrimSpace((*Answer)[0:ind])
		*AnswerParts = append(*AnswerParts, jsonStr)
		//prepare left data for next loop
		*Answer = (*Answer)[ind+3:]
		return true
	}

	//capture 3 of ```json{...}```
	JsonCompletedN := func(Answer string) (AnswerParts []string, done bool) {
		//chatgpt may return json without json tag
		//try to capture  `[...]` or `{...}`
		if !strings.Contains(Answer, "```json") {
			ind1, ind2 := strings.Index(Answer, "["), strings.LastIndex(Answer, "]")
			ind1_, ind2_ := strings.Index(Answer, "{"), strings.LastIndex(Answer, "}")
			// use outer [] or {}
			if ind1 >= 0 && ind2 > 0 && ind1 < ind2 {
				Answer = "```json" + Answer[ind1:ind2+1] + "```"
			} else if ind1_ >= 0 && ind2_ > 0 && ind1_ < ind2_ {
				Answer = "```json" + Answer[ind1_:ind2_+1] + "```"
			}
		}
		//中间的，缺少有效json标记情形的
		Answer = strings.Replace(Answer, "}\n\n{", "}\n``` \n```json{", -1)
		Answer = strings.Replace(Answer, "]\n\n[", "]\n``` \n```json[", -1)
		//capture 3 of ```json{...}```
		AnswerParts = []string{}
		for extractSucess := extractOneJson(&Answer, &AnswerParts); extractSucess; extractSucess = extractOneJson(&Answer, &AnswerParts) {
		}
		return AnswerParts, len(AnswerParts) == N
	}
	return JsonCompletedN
}
