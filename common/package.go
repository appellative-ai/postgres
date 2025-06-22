package common

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	LocationName = "location"
	PathName2    = "path"
	StatusName   = "status"
)

func LocationValues(h http.Header) (path string, status int, ok bool) {
	if h == nil {
		return
	}
	foundPath := false
	foundStatus := false
	for _, v := range h.Values(LocationName) {
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
