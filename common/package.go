package common

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	PostgresOverride = "x-postgres-override"
	PathName2        = "path"
	StatusName       = "status"
)

func QueryValues(h http.Header) (path string, status int, ok bool) {
	if h == nil {
		return
	}
	foundPath := false
	foundStatus := false
	for _, v := range h.Values(PostgresOverride) {
		t := strings.Split(v, "=")
		if len(t) != 2 {
			continue
		}
		if t[0] == PathName2 {
			path = t[1]
			foundPath = true
			continue
		}
		if t[0] == StatusName {
			status, _ = strconv.Atoi(t[1])
			foundStatus = true
			continue
		}
	}
	if foundPath || foundStatus {
		ok = true
		if status == 0 {
			status = http.StatusOK
		}
	}
	return
}
