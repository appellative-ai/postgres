package query

import (
	"github.com/behavioral-ai/postgres/common"
	"net/http"
	"strconv"
	"strings"
)

const (
	pathName   = "path"
	statusName = "status"
)

func queryValues(h http.Header) (path string, status int, ok bool) {
	if h == nil {
		return
	}
	pathFound := false
	statusFound := false
	for _, v := range h.Values(common.PostgresOverride) {
		t := strings.Split(v, "=")
		if len(t) != 2 {
			continue
		}
		if t[0] == pathName {
			path = t[1]
			pathFound = true
			continue
		}
		if t[0] == statusName {
			status, _ = strconv.Atoi(t[1])
			statusFound = true
			continue
		}
	}
	if pathFound || statusFound {
		ok = true
		if status == 0 {
			status = http.StatusOK
		}
	}
	return
}
