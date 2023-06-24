package client_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/0x00b/gobbq/components/gate/client"
	"github.com/0x00b/gobbq/conf"
	"github.com/0x00b/gobbq/engine/entity"
	"github.com/0x00b/gobbq/example/exampb"
	"github.com/0x00b/gobbq/tool/secure"
	"github.com/0x00b/gobbq/xlog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ClientService struct {
	entity.Entity
}

func (*ClientService) SayHello(c entity.Context, req *exampb.SayHelloRequest) (*exampb.SayHelloResponse, error) {

	xlog.Println("server req", req.String())

	return &exampb.SayHelloResponse{Text: fmt.Sprintf("response:%s", req.Text)}, nil
}

func TestWSClient(m *testing.T) {

	xlog.Init("trace", true, true, &lumberjack.Logger{
		Filename:  "./client.log",
		MaxAge:    7,
		LocalTime: true,
	}, xlog.DefaultLogTag{})
	conf.Init("client.yaml")

	client := client.NewClient(&exampb.ClientEntityDesc, &ClientService{})

	wg := sync.WaitGroup{}
	for i := 0; i < 1; i++ {
		wg.Add(1)

		i := i
		secure.GO(func() {

			// client := client.NewClient(&exampb.ClientEntityDesc, &ClientService{})

			es := exampb.NewEchoSvc2ServiceClient()
			cc, _ := client.Context().Copy()
			rsp, err := es.SayHello(cc, &exampb.SayHelloRequest{
				Text: fmt.Sprintf("%d", i),
				// Text:     "hello request",
				CLientID: uint64(client.EntityID()),
			})
			if err != nil {
				xlog.Errorln(err)
			}
			xlog.Infoln("rsp:", rsp)
			wg.Done()
		})
	}

	wg.Wait()

}
