package redis

import (
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/tool/sync2"
	"github.com/0x00b/gobbq/xlog"
	"github.com/go-redis/redis/v8"
)

var (
	// ConfigDSN config dsn description
	ConfigDSN = "./conf/app.yaml"

	// RedisSection redis
	RedisSection = "redis"

	// Decrypt func
	// Decrypt      config.DecryptFunc
	// configClient config.Client

	once sync2.OnceSucc
)

func initClient() error {
	var err error
	// configClient, err = config.New(&config.Option{
	// 	DSN: ConfigDSN,
	// })
	return err
}

// 1: log
// 2: trace todo
// 3: metrics  todo
type DefaultHook struct{}

// BeforeProcess 前回调
func (DefaultHook) BeforeProcess(ctx entity.Context, cmd redis.Cmder) (entity.Context, error) {
	xlog.Info(ctx, cmd.String())
	return ctx, nil
}

// AfterProcess 后回调
func (DefaultHook) AfterProcess(ctx entity.Context, cmd redis.Cmder) error {
	e := cmd.Err()
	if e != nil {
		xlog.Errorln(ctx, e)
	}
	return nil
}

// BeforeProcessPipeline 前pipeline
func (DefaultHook) BeforeProcessPipeline(ctx entity.Context, cmds []redis.Cmder) (entity.Context, error) {
	for _, cmd := range cmds {
		xlog.Info(ctx, cmd.String())
	}

	return ctx, nil
}

// AfterProcessPipeline 后pipeline
func (DefaultHook) AfterProcessPipeline(ctx entity.Context, cmds []redis.Cmder) error {
	for _, cmd := range cmds {
		e := cmd.Err()
		if e != nil {
			xlog.Errorln(ctx, e)
		}
	}
	return nil
}

func (DefaultHook) beforeAction(ctx entity.Context, cmd redis.Cmder) entity.Context {
	// log

	// metric
	return nil
}

func (DefaultHook) afterAction(ctx entity.Context, cmd redis.Cmder) {
	// log

	// metric
}

// NewClientFromConf init client from conf
// func NewClientFromConf(instanceName string) (*redis.Client, error) {
// 	err := once.Do(initClient)
// 	if err != nil {
// 		return nil, err
// 	}
// 	opt := &redis.Options{}
// 	key, err := config.ParseKey(ConfigDSN)
// 	if err != nil {
// 		return nil, err
// 	}
// 	err = configClient.Get(key, RedisSection).Unmarshal(opt)
// 	if err != nil {
// 		return nil, err
// 	}
// 	cli := redis.NewClient(opt)
// 	cli.AddHook(&DefaultHook{})
// 	return cli, nil
// }
