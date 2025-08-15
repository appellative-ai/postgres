package operations

import (
	"fmt"
	"net/http"
)

const (
	protocol       = "message-based"
	postgresScheme = "postgres"
)

// Request - contains data needed to build the SQL statement related to the uri
type requestT struct {
	method string
	uri    string
	h      http.Header
}

func newRequest(name, method string) *requestT {
	r := new(requestT)
	r.method = method
	r.uri = name
	r.h = make(http.Header)
	return r
}

func (r *requestT) Header() http.Header {
	return r.h
}

func (r *requestT) Method() string { return r.method }

func (r *requestT) Url() string {
	return r.uri
}

func (r *requestT) Protocol() string {
	return protocol
}

func buildUri(root, resource string) string {
	return fmt.Sprintf("%v://%v/%v:%v/%v/%v", postgresScheme, "host-name", "invalid-domain", "database-name", root, resource)
	//originUrn(nid, nss, test) //fmt.Sprintf("urn:%v.%v.%v:%v.%v", nid, o.Region, o.Zone, nss, test)
}
