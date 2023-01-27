package es

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/0x00b/gobbq/log"
	"github.com/elastic/go-elasticsearch/v8/estransport"
)

var _ estransport.Logger = Monitor{}

// log, metrics
type Monitor struct {
	DumpHttp bool
}

// LogRoundTrip should not modify the request or response, except for consuming and closing the body.
// Implementations have to check for nil values in request and response.
func (m Monitor) LogRoundTrip(req *http.Request, res *http.Response, err error, t time.Time, d time.Duration) error {
	c := req.Context()

	if m.DumpHttp {
		body, _ := httputil.DumpRequest(req, true)
		log.Infoln(c, "[Dump]:", string(body))
	}

	// if m.RequestBodyEnabled() && req != nil && req.Body != nil && req.Body != http.NoBody {
	// 	var buf bytes.Buffer
	// 	if req.GetBody != nil {
	// 		b, _ := req.GetBody()
	// 		buf.ReadFrom(b)
	// 	} else {
	// 		buf.ReadFrom(req.Body)
	// 	}
	// 	//out buf
	// }
	if m.ResponseBodyEnabled() && res != nil && res.Body != nil && res.Body != http.NoBody {
		defer res.Body.Close()
		buf, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Errorln(c, err)
		} else {
			var rsp map[interface{}]interface{}
			err := json.Unmarshal(buf, &rsp)
			if err != nil {
				log.Errorln(c, err)
			} else {
				log.Infoln(c, log.String(rsp))
			}
		}
	}
	if err != nil {
		log.Errorln(c, err)
	}
	return nil
}

// RequestBodyEnabled makes the client pass a copy of request body to the logger.
func (Monitor) RequestBodyEnabled() bool {
	return true
}

// ResponseBodyEnabled makes the client pass a copy of response body to the logger.
func (Monitor) ResponseBodyEnabled() bool {
	return true
}
