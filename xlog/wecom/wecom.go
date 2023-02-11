package wecom

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	key               = "key"
	defaultWebhookUrl = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key="
)

var (
	moreThanErrorLevels     = []logrus.Level{logrus.ErrorLevel, logrus.PanicLevel, logrus.FatalLevel}
	ErrWebHookKeyIsExpected = fmt.Errorf("webhook key is expected but not found")
)

type Options struct {
	WebhookUrl string
	WebhookKey string
	Levels     []logrus.Level
	Timeout    time.Duration
}

type LogrusHook struct {
	webhook string
	levels  []logrus.Level
	client  *http.Client
}

func NewLogrusHook(opts Options) (logrus.Hook, error) {
	if opts.WebhookKey == "" {
		return nil, ErrWebHookKeyIsExpected
	}
	if opts.WebhookUrl == "" {
		opts.WebhookUrl = defaultWebhookUrl
	}
	if len(opts.Levels) == 0 {
		opts.Levels = moreThanErrorLevels
	}

	u, err := url.Parse(opts.WebhookUrl)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Set(key, opts.WebhookKey)
	u.RawQuery = query.Encode()

	return &LogrusHook{webhook: u.String(), levels: moreThanErrorLevels, client: &http.Client{Timeout: opts.Timeout}}, nil
}

func (hook *LogrusHook) Levels() []logrus.Level {
	return hook.levels
}
