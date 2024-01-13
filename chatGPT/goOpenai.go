package chatGPT

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"strings"

// 	"github.com/rs/zerolog/log"
// 	openai "github.com/yangkequn/go-openai"
// 	"github.com/yangkequn/saavuu/api"
// 	"github.com/yangkequn/saavuu/data"
// )

// //rq := &InChatGPT{Text: Prompt, Model: openai.GPT4, AccessToken: "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", BaseUrl: FakeOpenAIUrl + "/v1"}

// var ApiChatGptAPI, _ = api.Api(func(param *InChatGPT) (Answer *OutChatGPT, err error) {
// 	var (
// 		resp                      openai.ChatCompletionResponse
// 		LastMessage, responseText string
// 		client                    *openai.Client

// 		gptChat = GptChatSession(param)
// 	)
// 	//set default model to  GPT3Dot5Turbo16K
// 	if param.Model == "" {
// 		param.Model = openai.GPT3Dot5Turbo16K
// 	}
// 	if param.AccessToken == "" {
// 		log.Fatal().AnErr("param", fmt.Errorf("AccessToken is empty")).Msg("ApiChatGptAPI")
// 	}

// 	if LastMessage, err = GptCacheGet(param.Model, param.Prompt); len(LastMessage) > 20 {
// 		return &OutChatGPT{Answer: LastMessage}, nil
// 	}

// 	if param.BaseUrl != "" {
// 		client = openai.NewClientWithConfig(openai.DefaultConfigWithBaseUrl(param.AccessToken, param.BaseUrl))
// 	} else {
// 		client = openai.NewClientWithConfig(openai.DefaultConfig(param.AccessToken))
// 	}

// 	var req = openai.ChatCompletionRequest{
// 		Model:    param.Model,
// 		Messages: []openai.ChatCompletionMessage{{Role: openai.ChatMessageRoleSystem, Content: "you are a helpful chatbot"}},
// 		Stream:   false,
// 	}

// 	//if not all messege is returned by ChatGPT, the continue to get the rest by messege "continue"
// 	//https://github.com/yangkequn/go-openai/blob/master/examples/chatbot/main.go
// 	//https://github.com/sashabaranov/go-openai/blob/master/examples/completion/main.go
// 	for i := 0; i < int(param.maxDialogueN); i++ {
// 		if i == 0 {
// 			req.Messages = append(req.Messages, openai.ChatCompletionMessage{
// 				Role:    openai.ChatMessageRoleUser,
// 				Content: param.Prompt,
// 			})
// 		} else {
// 			req.Messages = append(req.Messages, resp.Choices[0].Message)
// 			req.Messages = append(req.Messages, openai.ChatCompletionMessage{
// 				Role:    openai.ChatMessageRoleUser,
// 				Content: "继续回答",
// 			})
// 		}
// 		if param.BaseUrl != "" {
// 			if resp, err = client.CreateChatCompletionUseBaseUrl(context.Background(), req); err != nil {
// 				fmt.Printf("ChatCompletion error: %v\n", err)
// 				return
// 			}
// 		} else {
// 			if resp, err = client.CreateChatCompletion(context.Background(), req); err != nil {
// 				fmt.Printf("ChatCompletion error: %v\n", err)
// 				return
// 			}

// 		}
// 		Content := resp.Choices[0].Message.Content
// 		if i != 0 {
// 			//多段对话的衔接，去掉多余的markdown标记
// 			if _lines := strings.Split(Content, "```\n"); len(_lines) > 1 {
// 				Content = strings.Join(_lines[1:], "```\n")
// 			}
// 			if _lines := strings.Split(Content, "```json\n"); len(_lines) > 1 {
// 				Content = strings.Join(_lines[1:], "```json\n")
// 			}
// 			//if lastline  repeated in the response, then remove last line from responseText
// 			lines := strings.Split(responseText, "\n")
// 			if len(lines) > 1 && strings.Index(Content, lines[len(lines)-1]) == 0 {
// 				responseText = strings.Join(lines[:len(lines)-1], "\n")
// 			}
// 		}
// 		responseText += Content
// 		fmt.Println("\n\nGPTResponseText\n", responseText)
// 		lastChar := []rune(responseText)[len([]rune(responseText))-1]

// 		//doc received from ChatGPT is  complete,  break
// 		//? could not be the ending char. it is a question to be answered
// 		if resp.Choices[0].FinishReason == "stop" || (lastChar == '.' || lastChar == '!' || lastChar == ';' || lastChar == '。' || lastChar == '！' || lastChar == '\n') {
// 			break
// 		}
// 	}
// 	if responseText == "" {
// 		return nil, errors.New("no text")
// 	}
// 	//remove MatchTobeContinue
// 	Answer = &OutChatGPT{Answer: responseText,
// 		Model: resp.Model,
// 		MsgID: resp.ID,
// 	}
// 	GptCacheSet(req.Model, param.Prompt, responseText)
// 	return Answer, nil

// })

// type InOpenAICompletion struct {
// 	Prompt    string
// 	MaxTokens int64
// 	Model     string
// }

// var ApiOpenAICompletion, _ = api.Api(func(param *InOpenAICompletion) (Answer *OutChatGPT, err error) {
// 	var (
// 		resp                 openai.CompletionResponse
// 		LastMessage, Content string
// 	)
// 	if len(param.Model) == 0 {
// 		return nil, fmt.Errorf("no model in ApiOpenAICompletion")
// 	}

// 	if LastMessage, err = GptCacheGet(param.Model, param.Prompt); err == nil && len(LastMessage) > 20 {
// 		return &OutChatGPT{Answer: LastMessage}, nil
// 	}
// 	if param.MaxTokens <= 16 {
// 		return nil, errors.New("MaxTokens should be greater than 16")
// 	}

// 	var client = openai.NewClient("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")
// 	// GPT-3.5-turbo是一个更高效的模型，提供了与GPT-3类似的性能，但是以更低的价格和更快的响应速度。
// 	// 它特别适用于构建应用程序，如聊天机器人、虚拟助手、代码辅助等。
// 	// GPT-3.5-turbo提供了一个简化的API，使用openai.ChatCompletion.create()方法进行调用，可以处理一系列的消息以进行多轮对话。
// 	// GPT-3-Davinci-002
// 	// GPT-3-Davinci-002是GPT-3系列中的一个更大、更强大的模型。
// 	// 它在文本生成的质量、创造力和理解复杂提示方面通常优于GPT-3.5-turbo，但也更昂贵并且响应速度更慢。
// 	// GPT-3-Davinci-002通常使用openai.Completion.create()方法进行调用，更适合单轮的文本生成任务。
// 	req := openai.CompletionRequest{
// 		Model:     param.Model,
// 		MaxTokens: int(param.MaxTokens),
// 		Prompt:    param.Prompt,
// 	}

// 	//if not all messege is returned by ChatGPT, the continue to get the rest by messege "continue"
// 	//https://github.com/yangkequn/go-openai/blob/master/examples/chatbot/main.go
// 	//https://github.com/sashabaranov/go-openai/blob/master/examples/completion/main.go

// 	resp, err = client.CreateCompletion(context.Background(), req)
// 	if err != nil {
// 		fmt.Printf("Completion error: %v\n", err)
// 		return
// 	}

// 	if Content = resp.Choices[0].Text; len(Content) == 0 {
// 		return nil, errors.New("no text")
// 	}
// 	Answer = &OutChatGPT{Answer: Content,
// 		Model: resp.Model,
// 		MsgID: resp.ID,
// 	}
// 	GptCacheSet(param.Model, param.Prompt, Content)
// 	return Answer, nil

// })

// var ApiGptGoOpenAI, _ = api.Api(func(req *InChatGPT) (Answer *OutChatGPT, err error) {
// 	var (
// 		resp                      openai.CompletionResponse
// 		LastMessage, responseText string
// 	)
// 	// FakeOpenAIUrl+"/imitate/v1"
// 	if len(req.BaseUrl) == 0 {
// 		return nil, fmt.Errorf("no baseUrl")
// 	} else if len(req.AccessToken) == 0 {
// 		return nil, fmt.Errorf("no AccessToken")
// 	} else if len(req.Model) == 0 {
// 		return nil, fmt.Errorf("no model in ApiGptGoOpenAI")
// 	}

// 	if LastMessage, err = GptCacheGet(req.Model, req.Prompt); err == nil && len(LastMessage) > 20 {
// 		return &OutChatGPT{Answer: LastMessage}, nil
// 	}

// 	var client = openai.NewClientWithConfig(openai.DefaultConfigWithBaseUrl(req.AccessToken, req.BaseUrl))
// 	var chatReq = openai.CompletionRequest{
// 		Model:     openai.GPT3Dot5Turbo16K,
// 		MaxTokens: 5,
// 		Prompt:    req.Prompt,
// 		Stream:    false,
// 	}

// 	if resp, err = client.CreateCompletion(context.Background(), chatReq); err != nil {
// 		fmt.Printf("ChatCompletion error: %v\n", err)
// 		return
// 	}

// 	if responseText = resp.Choices[0].Text; responseText == "" {
// 		return nil, errors.New("no text")
// 	}
// 	Answer = &OutChatGPT{Answer: responseText,
// 		Model: resp.Model,
// 		MsgID: resp.ID,
// 	}
// 	GptCacheSet(req.Model, req.Prompt, responseText)
// 	return Answer, nil

// })

// var ApiGptCompletionAPI, _ = api.Api(func(param *InChatGPT) (Answer *OutChatGPT, err error) {
// 	var (
// 		resp                      openai.CompletionResponse
// 		LastMessage, responseText string
// 		client                    *openai.Client
// 	)
// 	//set default model to  GPT3Dot5Turbo16K
// 	if param.Model == "" {
// 		return nil, errors.New("model is empty")
// 	} else if param.AccessToken == "" {
// 		return nil, errors.New("AccessToken is empty")
// 	}

// 	if LastMessage, err = GptCacheGet(param.Model, param.Prompt); err == nil && len(LastMessage) > 20 {
// 		return &OutChatGPT{Answer: LastMessage}, nil
// 	}

// 	if param.BaseUrl != "" {
// 		client = openai.NewClientWithConfig(openai.DefaultConfigWithBaseUrl(param.AccessToken, param.BaseUrl))
// 	} else {
// 		client = openai.NewClientWithConfig(openai.DefaultConfig(param.AccessToken))
// 	}

// 	var req = openai.CompletionRequest{
// 		Model:     param.Model,
// 		Prompt:    param.Prompt,
// 		MaxTokens: int(param.MaxTokens),
// 	}

// 	//if not all messege is returned by ChatGPT, the continue to get the rest by messege "continue"
// 	//https://github.com/yangkequn/go-openai/blob/master/examples/chatbot/main.go
// 	//https://github.com/sashabaranov/go-openai/blob/master/examples/completion/main.go

// 	if resp, err = client.CreateCompletion(context.Background(), req); err != nil {
// 		fmt.Printf("ChatCompletion error: %v\n", err)
// 		return
// 	}
// 	responseText = resp.Choices[0].Text

// 	//remove MatchTobeContinue
// 	Answer = &OutChatGPT{Answer: responseText,
// 		Model: resp.Model,
// 		MsgID: resp.ID,
// 	}
// 	GptCacheSet(param.Model, param.Prompt, responseText)
// 	return Answer, nil

// })
