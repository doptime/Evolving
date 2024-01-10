package gpt

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cespare/xxhash/v2"
	"github.com/elliotchance/pie/v2"
	"github.com/rs/zerolog/log"
	"github.com/tidwall/gjson"

	"github.com/danil/randuuid"
)

type InChatGPT struct {
	Prompt             string
	ParentMessegeID    string
	ConversationID     string
	Model              string
	AccessToken        string
	BaseUrl            string
	maxPackN           int64
	MaxTokens          int64
	RemoveConversation bool
	XYHelperAuthKey    string
}
type OutChatGPT struct {
	Answer         string
	MsgID          string
	ConversationID string
	Model          string
	UseCache       bool
}

func ApiGptChatOnce(req *InChatGPT) (Answer *OutChatGPT, err error) {
	var (
		message_id, LastMessage, answer string
		lines, lines1                   []string
	)
	Answer = &OutChatGPT{}
	type FakeAIDataAuthor struct {
		Role string `json:"role"`
	}
	type FakeAIDataContent struct {
		Content_type string   `json:"content_type"`
		Parts        []string `json:"parts"`
	}
	type FakeAIChatMessage struct {
		Id      string            `json:"id"`
		Author  FakeAIDataAuthor  `json:"author"`
		Content FakeAIDataContent `json:"content"`
	}
	type FakeAIChatData struct {
		Action            string               `json:"action"`
		Messages          []*FakeAIChatMessage `json:"messages"`
		Model             string               `json:"model"`
		Parent_message_id string               `json:"parent_message_id"`
		//Stream            bool                 `json:"stream"` 官方不支持

		//json skip empty
		Conversation_id string `json:"conversation_id,omitempty"`
		//"history _and training disabled": false
		HistoryAndTrainingDisabled bool `json:"history_and_training_disabled"`
	}
	//Conversation_id   string               `json:"conversation_id"`

	// prompt 提问的内容。
	// model 对话使用的模型，通常整个会话中保持不变。如 gpt-4-mobile, gpt-4, gpt-3.5
	// message_id 消息ID，通常使用str(uuid.uuid4())来生成一个。
	// parent_message_id 父消息ID。获取上一条回复的消息ID即可。
	// conversation_id 首次对话可不传。ChatGPT回复时可获取。
	// stream 是否使用流的方式输出内容，默认为：True.返回多条信息，最后一条是[DONE]。如果为false，只返回一条信息，不包含[DONE]
	//https://github.com/pengzhile/pandora/blob/master/doc/HTTP-API.md

	//自建反代理 https://pypi.org/project/nonebot-plugin-talk-with-chatgpt/

	hash := xxhash.Sum64String(req.Prompt)
	//prevent repeated uuid request
	uuid, _ := randuuid.New(int64(hash) ^ time.Now().UnixNano())
	message_id = uuid.String()
	//shoultd use jason.Marshal, rather than fmt.Sprintf, to avoid special char
	//https://github.com/pengzhile/pandora/blob/master/src/pandora/openai/api.py#L223
	// 官方api说明 https://platform.openai.com/docs/guides/gpt/chat-completions-api
	author, content := FakeAIDataAuthor{Role: "user"}, FakeAIDataContent{Content_type: "text", Parts: []string{req.Prompt}}
	messge := &FakeAIChatMessage{Id: message_id, Author: author, Content: content}
	postData := &FakeAIChatData{Action: "next", Messages: []*FakeAIChatMessage{messge}, Model: req.Model, Parent_message_id: req.ParentMessegeID,
		Conversation_id: req.ConversationID, HistoryAndTrainingDisabled: false}
	// Create a new request
	if answer, err = FakeAIPost(req.BaseUrl, req.AccessToken, req.XYHelperAuthKey, postData); err != nil {
		return nil, err
	} else if answer == "" {
		return nil, errors.New("no answer from chatGPT")
	}

	//support stream with multiple messege and with [DONE] at the end
	//take the last message as LastMessage
	if lines = strings.Split(answer, "\n"); len(lines) == 0 {
		return nil, errors.New("no answer from chatGPT")
	}
	lines1 = pie.Filter(lines, func(i string) bool {
		return strings.Contains(i, "finished_successfully\",")
	})
	if len(lines1) == 0 {
		if len(lines) > 0 {
			return nil, errors.New(lines[len(lines)-1])
		}
		return nil, errors.New("no finished_successfully")
	}
	lines = lines1

	//take the last message
	if LastMessage = lines[len(lines)-1]; strings.Contains(LastMessage, "data:") {
		LastMessage = LastMessage[5:]
	}
	if json := gjson.Parse(LastMessage); json.Get("messege").String() != "" {
		LastMessage = json.String()
	}
	//convert LastMessage to raw unicode string ,rather than encoded unicode string
	var obj interface{}
	if err = json.Unmarshal([]byte(LastMessage), &obj); err != nil {
		return nil, err
	}
	if bytes, err := json.Marshal(obj); err == nil {
		LastMessage = string(bytes)

	}
	err = DecodeFromData(LastMessage, Answer)
	return Answer, err
}
func DecodeFromData(Messege string, Answer *OutChatGPT) (err error) {
	js := gjson.Parse(Messege)
	messege, detail := js.Get("message"), js.Get("detail")
	//gpt plus restrictions
	if detailStr := detail.String(); detailStr != "" {
		return errors.New("ApiGptChatOnce: " + detailStr)
	}
	//fmt.Println(LastMessage)
	ConversationID, MsgID, Model := js.Get("conversation_id").String(), messege.Get("id").String(), messege.Get("metadata").Get("model_slug").String()
	if len(ConversationID) > 0 && len(MsgID) > 0 {
		Answer.ConversationID = ConversationID
		Answer.MsgID = MsgID
		Answer.Model = Model
	} else {
		return errors.New("ApiGptChatOnce: ConversationID or MsgID . RawString: " + Messege)
	}
	if len(Model) == 0 {
		log.Warn().Msgf("ApiGptChatOnce: Model is empty. RawString: %s", Messege)
	}
	//here parts is an array,concat all the parts to string to get the answer
	MsgMissing := fmt.Errorf("ApiGptChatOnce: messege content missing")
	for _, part := range messege.Get("content").Get("parts").Array() {
		Answer.Answer += part.String()
		MsgMissing = nil
	}
	return MsgMissing
}
func funcGptRaw(req *InChatGPT, terminated func(answer string) (result []string, done bool)) (Answer *OutChatGPT, err error) {
	var (
		LastMessage string
	)
	if req.BaseUrl == "" {
		return nil, errors.New("ChatGptRaw:BaseUrl is empty")
	}
	if req.Model == "" {
		return nil, errors.New("ChatGptRaw:Model is empty")
	}

	if Answer, err = LoadCache(req, LastMessage); err == nil {
		return Answer, nil
	}
	for i := 0; i < int(req.maxPackN); i++ {
		Answer, err = ApiGptChatOnce(req)
		if err != nil {
			return nil, err
		} else if Answer == nil || Answer.Answer == "" {
			return nil, errors.New("no answer from chatGPT")
		}
		//多段对话的衔接，去掉多余的markdown标记
		if i > 0 {
			RemovedHead := false
			for _, sep := range []string{"```\n", "```json\n", "继续回答：\n\n"} {
				if _lines := strings.Split(Answer.Answer, sep); len(_lines) > 1 && len(_lines[0]) < 60 {
					//log.Warn().Msgf("before RemovedHead: %s", Answer.Answer)
					Answer.Answer = strings.Join(_lines[1:], sep)
					RemovedHead = true
					break
				}
			}
			if !RemovedHead {
				//log.Warn().Msgf("funcGptRaw: no markdown head removed: %s", Answer.Answer)
			}
		}
		//if last line of last answer repeated in new answer, then dispose the last line in last answer
		if llines := strings.Split(LastMessage, "\n"); i > 0 && len(llines) > 0 && len(llines[len(llines)-1]) > 0 {
			lstLine := llines[len(llines)-1]
			ind := strings.Index(Answer.Answer, lstLine)
			if ind >= 0 && ind < 5 {
				LastMessage = strings.Join(llines[:len(llines)-1], "\n")
			}
		}
		LastMessage = LastMessage + Answer.Answer
		if terminated == nil {
			break
		} else if _, done := terminated(LastMessage); done {
			break
		}
		req.Prompt = "继续回答"
		req.ParentMessegeID = Answer.MsgID
		req.ConversationID = Answer.ConversationID
	}
	//raise error if no answer
	if LastMessage == "" {
		return nil, errors.New("no answer from chatGPT")
	}
	GptCacheSet(req.Model, req.Prompt, LastMessage)
	//remove the chatgpt conversation
	if Answer.ConversationID != "" && req.RemoveConversation {
		DeleteUrl := fmt.Sprintf(req.BaseUrl+"/%s", Answer.ConversationID)
		FakeAIDel(DeleteUrl, req.AccessToken)
	}
	Answer.Answer = LastMessage
	return Answer, nil
}
