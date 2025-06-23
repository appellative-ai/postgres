package retrieval

import (
	"fmt"
	"net/http"
)

const (
	protocol       = "message-based"
	postgresScheme = "postgres"
	queryRoot      = "retrieval"
	queryRouteName = "postgresql-retrieval"
	selectMethod   = "select"
)

type Attr struct {
	Key string
	Val any
}

// Request - contains data needed to build the SQL statement related to the uri
type request struct {
	resource  string
	template  string
	uri       string
	routeName string

	values  [][]any
	values2 map[string][]string
	attrs   []Attr
	where   []Attr
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

// BuildWhere - build the []Attr based on the URL retrieval parameters
func buildWhere(values map[string][]string) []Attr {
	if len(values) == 0 {
		return nil
	}
	var where []Attr
	for k, v := range values {
		where = append(where, Attr{Key: k, Val: v[0]})
	}
	return where
}
