package exec

import (
	"fmt"
	"github.com/behavioral-ai/postgres/common"
	"net/http"
	"time"
)

const (
	protocol       = "message-based"
	postgresScheme = "postgres"
	queryRoot      = "retrieval"
	execRoot       = "exec"
	pingRoot       = "ping"

	selectMethod = "select"
	insertMethod = "insert"
	updateMethod = "update"
	deleteMethod = "delete"
	insertCmd    = 1
	updateCmd    = 2
	deleteCmd    = 3

	nullExpectedCount = int64(-1)

	PkgPath       = "github/behavioral-ai/postgres/pgxsql"
	userConfigKey = "user"
	pswdConfigKey = "pswd"
	uriConfigKey  = "uri"

	insertRouteName = "postgres-insert"
	updateRouteName = "postgres-update"
	deleteRouteName = "postgres-delete"
)

// Request - contains data needed to build the SQL statement related to the uri
type request struct {
	expectedCount int64
	cmd           int
	duration      time.Duration

	resource  string
	template  string
	uri       string
	routeName string

	values  [][]any
	values2 map[string][]string
	attrs   []common.Attr
	where   []common.Attr
	args    []any
	error   error
	h       http.Header
}

func newRequest(cmd int, resource, template, uri, routeName string) *request {
	r := new(request)
	r.expectedCount = nullExpectedCount
	r.cmd = cmd

	r.resource = resource
	r.template = template
	r.uri = uri
	r.routeName = routeName
	r.duration = -1
	r.h = make(http.Header)
	return r
}

func (r *request) Header() http.Header {
	return r.h
}

func (r *request) Method() string {
	switch r.cmd {
	case insertCmd:
		return insertMethod
	case updateCmd:
		return updateMethod
	case deleteCmd:
		return deleteMethod
	}
	return "unknown"
}

func (r *request) Url() string {
	return r.uri
}

func (r *request) Protocol() string {
	return protocol
}

func buildUri(root, resource string) string {
	return fmt.Sprintf("%v://%v/%v:%v/%v/%v", postgresScheme, "host-name", "invalid-domain", "database-name", root, resource)
	//originUrn(nid, nss, test) //fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, o.Region, o.Zone, nss, test)
}

func newInsertRequest(resource, template string, values [][]any, args ...any) *request {
	r := newRequest(insertCmd, resource, template, buildUri(execRoot, resource), insertRouteName)
	r.values = values
	r.args = args
	return r
}

func newUpdateRequest(resource, template string, attrs []common.Attr, where []common.Attr, args ...any) *request {
	r := newRequest(updateCmd, resource, template, buildUri(execRoot, resource), updateRouteName)
	r.attrs = attrs
	r.where = where
	r.args = args
	return r
}

func newDeleteRequest(resource, template string, where []common.Attr, args ...any) *request {
	r := newRequest(deleteCmd, resource, template, buildUri(execRoot, resource), deleteRouteName)
	r.where = where
	r.args = args
	return r
}
