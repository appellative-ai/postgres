package query

import (
	"fmt"
	"github.com/behavioral-ai/postgres/common"
	"net/http"
)

const (
	protocol       = "message-based"
	postgresScheme = "postgres"
	queryRoot      = "query"
	queryRouteName = "postgresql-query"

	//execRoot       = "exec"
	//pingRoot       = "ping"

	selectMethod = "select"
	//insertMethod = "insert"
	//updateMethod = "update"
	//deleteMethod = "delete"
	//pingMethod   = "ping"

	//selectCmd = 0
	//insertCmd = 1
	//updateCmd = 2
	//deleteCmd = 3
	//pingCmd   = 4

	//nullExpectedCount = int64(-1)
)

// Request - contains data needed to build the SQL statement related to the uri
type request struct {
	//expectedCount int64
	//cmd           int
	//duration      time.Duration

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

func newRequest(resource, template string) *request {
	r := new(request)
	r.resource = resource
	r.template = template
	r.uri = buildUri(queryRoot, resource)
	r.routeName = queryRouteName
	r.h = make(http.Header)
	return r
}

func (r *request) Header() http.Header {
	return r.h
}

func (r *request) Method() string {
	return selectMethod
}

func (r *request) Url() string {
	return r.uri
}

func (r *request) Protocol() string {
	return protocol
}

func buildUri(root, resource string) string {
	return fmt.Sprintf("%v://%v/%v:%v/%v/%v", postgresScheme, "host-name", "invalid-domain", "database-name", root, resource)
}

// buildQueryUri - build an uri with the Query NSS
func buildQueryUri(resource string) string {
	return buildUri(queryRoot, resource)
}

// BuildWhere - build the []Attr based on the URL query parameters
func buildWhere(values map[string][]string) []common.Attr {
	if len(values) == 0 {
		return nil
	}
	var where []common.Attr
	for k, v := range values {
		where = append(where, common.Attr{Key: k, Val: v[0]})
	}
	return where
}
