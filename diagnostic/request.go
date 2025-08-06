package diagnostic

import (
	"fmt"
	"net/http"
)

const (
	protocol       = "message-based"
	postgresScheme = "postgres"
)

// Request - contains data needed to build the SQL statement related to the uri
type request struct {
	method string
	uri    string
	h      http.Header
}

func newRequest(name, method string) *request {
	r := new(request)
	r.method = method
	r.uri = name
	r.h = make(http.Header)
	return r
}

func (r *request) Header() http.Header {
	return r.h
}

func (r *request) Method() string { return r.method }

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
