package apiChatGPT

import "strings"

// GptResponseParseToJS_MD_Bash handles both JSON and Bash script sections.
type GPTResponseCode struct {
	Type string
	Text string
}

// extractOneSection extracts one section (either JSON or Bash) from the answer.
func extractOneSection(Answer *string, AnswerParts *[]*GPTResponseCode) (success bool) {
	var (
		ind int
	)
	if ind = strings.Index(*Answer, "```"); ind < 0 {
		return false
	}
	var sectionTag string = ""
	for _, _char := range (*Answer)[ind+3:] {
		if _char == '\n' {
			break
		}
		sectionTag += string(_char)
	}
	if len(sectionTag) == 0 {
		return false
	}

	*Answer = (*Answer)[ind+len(sectionTag)+3:]
	if ind = strings.Index(*Answer, "```"); ind < 0 {
		return false
	}
	sectionStr := strings.TrimSpace((*Answer)[:ind])
	*AnswerParts = append(*AnswerParts, &GPTResponseCode{Type: sectionTag, Text: sectionStr})
	*Answer = (*Answer)[ind+3:]
	return true
}

// CompletedNSections extracts N sections of the specified type from the answer.
// json,markdown,bash,html,latex,python,ruby,sql,swift,typescript,css,go,java,javascript,php,scala,xml,...
func GptResponseParseToCodes(Answer string) (gptResponseSnippets []*GPTResponseCode) {
	//chatgpt may return json without json tag
	//try to capture  `[...]` or `{...}`
	if !strings.Contains(Answer, "```") {
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

	// Capture N sections of the specified type.
	gptResponseSnippets = []*GPTResponseCode{}
	for extractSucess := true; extractSucess; {
		extractSucess = extractOneSection(&Answer, &gptResponseSnippets)
	}
	return gptResponseSnippets
}
func GptResponseParseToJS_MD_Bash1(N int) func(Answer string) (JSParts []string, done bool) {
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
