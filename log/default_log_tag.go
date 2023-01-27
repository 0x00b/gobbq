package log

import "context"

type DefaultLogTag struct {
}

func (DefaultLogTag) GetTags(c context.Context) map[string]string {
	return map[string]string{
		// "TraceId": trace.GetTraceID(c),
		// "Action":  ago.Meta(c).Action,
	}
}
