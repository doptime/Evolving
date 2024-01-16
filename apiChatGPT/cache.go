package apiChatGPT

import (
	"errors"
	"fmt"
	"strings"

	"github.com/yangkequn/saavuu/data"
)

func GptCacheGet(model string, Prompt string) (LastMessage string, err error) {
	if len(model) == 0 {
		return "", errors.New("GptCacheSet: model is empty")
	} else if len(Prompt) == 0 {
		return "", errors.New("GptCacheSet: Prompt is empty")
	}
	var keyGptCache = data.New[string, string](model)
	return keyGptCache.Get(Prompt)
}
func GptCacheSet(model string, Prompt string, LastMessage string) (err error) {
	if len(model) == 0 {
		return errors.New("GptCacheSet: model is empty")
	} else if len(Prompt) == 0 {
		return errors.New("GptCacheSet: Prompt is empty")
	} else if len(LastMessage) == 0 {
		return errors.New("GptCacheSet: LastMessage is empty")
	}

	var keyGptCache = data.New[string, string](model)
	return keyGptCache.Set(Prompt, LastMessage, 0)
}
func GptCacheRemove(model string, Prompt string) (err error) {
	if len(model) == 0 {
		return errors.New("GptCacheSet: model is empty")
	} else if len(Prompt) == 0 {
		return errors.New("GptCacheSet: Prompt is empty")
	}
	var keyGptCache = data.New[string, string](model)
	return keyGptCache.Del(Prompt)
}
func (ss *ChatGPTSession) LoadCache(Prompt string) (Answer *OutChatGPT, err error) {
	var (
		LastMessage string
	)
	if ss.Model == "" {
		return nil, errors.New("LoadCache:Model is empty")
	}
	LastMessage, err = GptCacheGet(ss.Model, Prompt)
	if err != nil || len(LastMessage) < 20 {
		return nil, fmt.Errorf("LoadCache: noCache")
	}
	Answer = &OutChatGPT{UseCache: true}
	if strings.Contains(LastMessage, "conversation_id") {
		Answer.Answer, err = ss.DecodeFromResponseData(LastMessage)
	} else {
		Answer.Answer = LastMessage
	}
	return Answer, err
}
