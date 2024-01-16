package apiChatGPT

import (
	"time"

	"github.com/rs/zerolog/log"
)

var Models = map[string]string{"gpt-3.5-turbo-1106": "gpt-3.5-turbo-1106", "gpt-4-plugins": "gpt-4-plugins", "gpt-4-gizmo": "gpt-4-gizmo"}

var NewChatSession = func(Model string, maxPackN int64) *ChatGPTSession {
	if len(Model) == 0 {
		log.Fatal().Msg("no Model specified in ApiChatGptXY")
	}
	rq := &ChatGPTSession{MaxTokens: 8192, Model: Model, BaseUrl: "https://demo.xyhelper.cn/backend-api/conversation",
		RemoveConversation: true, maxDialogueN: maxPackN, XYHelperAuthKey: "xyhelper"}

	if rq.BaseUrl == "" {
		log.Fatal().Msg("ChatGptRaw:BaseUrl is empty")
	}
	if rq.Model == "" {
		log.Fatal().Msg("ChatGptRaw:Model is empty")
	}

	//fill access token
	for rq.AccessToken = TakeOneValidAccessToken(rq.Model); len(rq.AccessToken) == 0; rq.AccessToken = TakeOneValidAccessToken(rq.Model) {
		log.Error().Msg("no AccessToken available")
		time.Sleep(time.Second * 60)
	}
	//remove the chatgpt conversation
	defer rq.RemoveChatHistory()

	return rq

}
