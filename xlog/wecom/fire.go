package wecom

import (
	"bytes"
	"encoding/json"

	"github.com/sirupsen/logrus"
)

const (
	msgTypeText     = "text"
	contentTypeJSON = "application/json;charset=utf-8"
)

type message struct {
	MsgType string     `json:"msgtype"`
	Text    msgContent `json:"text"`
}

type msgContent struct {
	Content       string   `json:"content"`
	MentionedList []string `json:"mentioned_list,omitempty"`
}

func (hook *LogrusHook) Fire(entry *logrus.Entry) error {
	go func() {
		_ = hook.fire(entry)
	}()
	return nil
}

func (hook *LogrusHook) fire(entry *logrus.Entry) error {
	defer func() { _ = recover() }()

	line, err := entry.String()
	if err != nil {
		return err
	}

	msg := &message{
		MsgType: msgTypeText,
		Text: msgContent{
			Content: line,
		},
	}

	byts, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := hook.client.Post(hook.webhook, contentTypeJSON, bytes.NewBuffer(byts))
	if resp != nil {
		_ = resp.Body.Close()
	}
	return err
}
