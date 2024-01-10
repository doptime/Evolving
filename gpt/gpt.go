package gpt

import (
	"context"
	"errors"
	"fmt"

	"github.com/samber/lo"

	openai "github.com/yangkequn/go-openai"
	"github.com/yangkequn/saavuu/api"
)

var Models = map[string]string{"3.5": "gpt-3.5-turbo-1106", "4p": "gpt-4-plugins", "4g": "gpt-4-gizmo"}

var ApiChatGptXY = func(Prompt string, Model string, maxPackN int64, terminated func(answer string) (result []string, done bool)) (Answer *OutChatGPT, err error) {
	if len(Model) == 0 {
		return &OutChatGPT{Answer: ""}, errors.New("no Model specified in ApiChatGptXY")
	}
	rq := &InChatGPT{MaxTokens: 4096, Prompt: Prompt, Model: Model, BaseUrl: "https://demo.xyhelper.cn/backend-api/conversation",
		RemoveConversation: true, maxPackN: maxPackN, XYHelperAuthKey: "xyhelper"}
	if rq.AccessToken = QuotaCheckOK(Model); len(rq.AccessToken) == 0 {
		return &OutChatGPT{Answer: ""}, errors.New("no AccessToken available")
	}
	Answer, err = funcGptRaw(rq, terminated)

	ReportError("GptLi", rq.AccessToken, Model, err)
	if err != nil {
		AccountInfo := lo.Filter(AuthTokenSlice, func(v *AuthToken, index int) bool {
			return v.AccessToken == rq.AccessToken
		})[0].AccountInfo
		err = fmt.Errorf("OrderID:%v ,%s", AccountInfo, err.Error())
	}
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

var ApiGptCompletionAPI, _ = api.Api(func(param *InChatGPT) (Answer *OutChatGPT, err error) {
	var (
		resp                      openai.CompletionResponse
		LastMessage, responseText string
		client                    *openai.Client
	)
	//set default model to  GPT3Dot5Turbo16K
	if param.Model == "" {
		return nil, errors.New("model is empty")
	} else if param.AccessToken == "" {
		return nil, errors.New("AccessToken is empty")
	}

	if LastMessage, err = GptCacheGet(param.Model, param.Prompt); err == nil && len(LastMessage) > 20 {
		return &OutChatGPT{Answer: LastMessage}, nil
	}

	if param.BaseUrl != "" {
		client = openai.NewClientWithConfig(openai.DefaultConfigWithBaseUrl(param.AccessToken, param.BaseUrl))
	} else {
		client = openai.NewClientWithConfig(openai.DefaultConfig(param.AccessToken))
	}

	var req = openai.CompletionRequest{
		Model:     param.Model,
		Prompt:    param.Prompt,
		MaxTokens: int(param.MaxTokens),
	}

	//if not all messege is returned by ChatGPT, the continue to get the rest by messege "continue"
	//https://github.com/yangkequn/go-openai/blob/master/examples/chatbot/main.go
	//https://github.com/sashabaranov/go-openai/blob/master/examples/completion/main.go

	if resp, err = client.CreateCompletion(context.Background(), req); err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}
	responseText = resp.Choices[0].Text

	//remove MatchTobeContinue
	Answer = &OutChatGPT{Answer: responseText,
		Model: resp.Model,
		MsgID: resp.ID,
	}
	GptCacheSet(param.Model, param.Prompt, responseText)
	return Answer, nil

})
