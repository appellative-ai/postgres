package exec

import (
	"github.com/behavioral-ai/postgres/common"
	"net/http"
	"strconv"
	"strings"
)

const (
	countName  = "count"
	statusName = "status"
)

func execValues(h http.Header) (count, status int, ok bool) {
	if h == nil {
		return
	}
	countFound := false
	statusFound := false
	for _, v := range h.Values(common.PostgresOverride) {
		t := strings.Split(v, "=")
		if len(t) != 2 {
			continue
		}
		if t[0] == countName {
			count, _ = strconv.Atoi(t[1])
			countFound = true
			continue
		}
		if t[0] == statusName {
			status, _ = strconv.Atoi(t[1])
			statusFound = true
			continue
		}
	}
	if countFound || statusFound {
		ok = true
		if status == 0 {
			status = http.StatusOK
		}
	}
	return
}
