package pgxsql

import (
	"github.com/behavioral-ai/core/access"
	"github.com/behavioral-ai/core/core"
	"net/http"
	"time"
)

func log(start time.Time, h http.Header, req *request, status *core.Status) {
	from := ""
	if h != nil {
		from = h.Get(core.XFrom)
	}
	cc := ""
	// TODO : determine if status caused by a timeout
	if status != nil {

	}
	access.Log(access.EgressTraffic, start, time.Since(start), req, status, access.Routing{From: from, Route: req.routeName, To: ""}, access.Controller{Timeout: req.duration, RateLimit: -1, RateBurst: -1, Code: cc})
}
