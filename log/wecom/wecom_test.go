package wecom

import (
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogrusHook(t *testing.T) {
	webhook, exists := os.LookupEnv("WEBHOOK_KEY")
	if !exists || webhook == "" {
		t.Skip("env WEBHOOK_KEY not found, skip")
	}

	hook, err := NewLogrusHook(Options{WebhookKey: webhook})
	assert.NoError(t, err)

	log := logrus.New()
	log.Hooks.Add(hook)

	log.Infoln("just print no, notify")
	log.Errorln("error message fire to wecom")

	// 异步发送，等待一段时间
	time.Sleep(time.Second)
}
