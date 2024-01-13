package apiChatGPT

import (
	"errors"
	"fmt"

	"github.com/tidwall/gjson"
)

func DecodeFromResponseData(Messege string, ConversationID, MsgID, ModelSlug *string) (Answer string, err error) {
	var (
		cid_exist, msgID_exist, model_slug_exist bool
		js                                       = gjson.Parse(Messege)
	)
	messege, detail := js.Get("message"), js.Get("detail")
	//gpt plus restrictions
	if detailStr := detail.String(); detailStr != "" {
		return "", errors.New(detailStr)
	}
	cid, msgID, model_slug := js.Get("conversation_id"), messege.Get("id"), messege.Get("metadata").Get("model_slug")
	if cid_exist = cid.Exists(); cid_exist {
		*ConversationID = cid.String()
	} else {
		return "", errors.New("gpt chat rsb missing conversation_id")
	}
	if msgID_exist = msgID.Exists(); msgID_exist {
		*MsgID = msgID.String()
	} else {
		return "", errors.New("gpt chat rsb missing message_id")
	}
	if model_slug_exist = model_slug.Exists(); model_slug_exist {
		*ModelSlug = model_slug.String()
	} else {
		return "", errors.New("gpt chat rsb missing model_slug")
	}
	//here parts is an array,concat all the parts to string to get the answer
	for _, part := range messege.Get("content").Get("parts").Array() {
		Answer += part.String()
	}
	if len(Answer) == 0 {
		return "", fmt.Errorf("ApiGptChatOnce: messege content missing")
	}
	return Answer, nil
}
