package apiChatGPT

import (
	"errors"

	"github.com/rs/zerolog/log"
)

var Models = map[string]string{"gpt-3.5-turbo-1106": "gpt-3.5-turbo-1106", "gpt-4-plugins": "gpt-4-plugins", "gpt-4-gizmo": "gpt-4-gizmo"}

var ApiChatGptXY = func(Model string, maxPackN int64, terminated func(answer string) (result []string, done bool)) func(Prompt string) (Answer *OutChatGPT, err error) {
	if len(Model) == 0 {
		log.Fatal().Msg("no Model specified in ApiChatGptXY")
	}
	rq := &InChatGPT{MaxTokens: 4096, Model: Model, BaseUrl: "https://demo.xyhelper.cn/backend-api/conversation",
		RemoveConversation: true, maxDialogueN: maxPackN, XYHelperAuthKey: "xyhelper"}
	for rq.AccessToken = QuotaCheckOK(Model); len(rq.AccessToken) == 0; {
		log.Error().Msg("no AccessToken available")
	}
	return GptChatSession(rq)

}
var ApiChatGptXYLongChat = func(Prompt string, Model string, maxPackN int64, terminated func(answer string) (result []string, done bool)) (Answer *OutChatGPT, err error) {
	if len(Model) == 0 {
		return &OutChatGPT{Answer: ""}, errors.New("no Model specified in ApiChatGptXY")
	}
	rq := &InChatGPT{MaxTokens: 4096, Model: Model, BaseUrl: "https://demo.xyhelper.cn/backend-api/conversation",
		RemoveConversation: true, maxDialogueN: maxPackN, XYHelperAuthKey: "xyhelper"}
	if rq.AccessToken = QuotaCheckOK(Model); len(rq.AccessToken) == 0 {
		return &OutChatGPT{Answer: ""}, errors.New("no AccessToken available")
	}
	Answer, err = GptLongDialogues(rq, Prompt)

	return Answer, err
}

// var ApiGPTCompletion35TurboInstruct, _ = api.ApiNamed("ChatGPCompletion35", func(Prompt string) (Answer *OutChatGPT, err error) {
// 	//rq := &InChatGPT{Text: Prompt, Model: openai.GPT3Dot5Turbo16K, BaseUrl: FakeOpenAIUrl + "/imitate/v1"}
// 	//rq := &InChatGPT{Text: Prompt, Model: openai.GPT3Dot5Turbo16K, AccessToken: "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"}
// 	//3.5 turbo key ,有限速
// 	//rq := &InChatGPT{Text: Prompt, Model: openai.GPT3Dot5TurboInstruct, AccessToken: "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"}
// 	//4.0 key ,无限速
// 	rq := &InChatGPT{MaxTokens: 4096, Prompt: Prompt, Model: openai.GPT3Dot5TurboInstruct, AccessToken: "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"}

// 	return ApiGptCompletionAPI(rq)
// })
// var ApiGPTCompletion40Turbo, _ = api.ApiNamed("ChatGPCompletion40Turbo", func(Prompt string) (Answer *OutChatGPT, err error) {
// 	rq := &InChatGPT{MaxTokens: 32768, Prompt: Prompt, Model: "gpt-4-1106-preview", AccessToken: "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"}

// 	return ApiChatGptAPI(rq)
// })
// var ApiGPTCompletion40, _ = api.ApiNamed("ChatGPCompletion40", func(Prompt string) (Answer *OutChatGPT, err error) {
// 	rq := &InChatGPT{MaxTokens: 32768, Prompt: Prompt, Model: openai.GPT40613, AccessToken: "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"}

// 	return ApiChatGptAPI(rq)
// })
