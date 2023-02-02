package log

import "github.com/0x00b/gobbq/engine/entity"

type DefaultLogTag struct {
}

func (DefaultLogTag) GetTags(c *entity.Context) map[string]string {
	return map[string]string{
		// "TraceId": trace.GetTraceID(c),
		// "Action":  ago.Meta(c).Action,
	}
}
