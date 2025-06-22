package query

import (
	"fmt"
	"github.com/behavioral-ai/postgres/common"
	"net/http"
)

func ExampleQueryValues() {
	var h http.Header

	path, status, ok := queryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Set(common.PostgresOverride, "test")
	path, status, ok = queryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h.Set(common.PostgresOverride, pathName)
	path, status, ok = queryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Add(common.PostgresOverride, pathName+"="+"file://test/data")
	path, status, ok = queryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Add(common.PostgresOverride, statusName+"="+"500")
	path, status, ok = queryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Add(common.PostgresOverride, pathName+"="+"file://test/data")
	h.Add(common.PostgresOverride, statusName+"="+"419")
	path, status, ok = queryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	//Output:
	//test: QueryValues() -> [path:] [status:0] [ok:false]
	//test: QueryValues() -> [path:] [status:0] [ok:false]
	//test: QueryValues() -> [path:] [status:0] [ok:false]
	//test: QueryValues() -> [path:file://test/data] [status:200] [ok:true]
	//test: QueryValues() -> [path:] [status:500] [ok:true]
	//test: QueryValues() -> [path:file://test/data] [status:419] [ok:true]

}
