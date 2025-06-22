package common

import (
	"fmt"
	"net/http"
)

func ExampleQueryValues() {
	var h http.Header

	path, status, ok := QueryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Set(PostgresOverride, "test")
	path, status, ok = QueryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h.Set(PostgresOverride, PathName2)
	path, status, ok = QueryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Add(PostgresOverride, PathName2+"="+"file://test/data")
	path, status, ok = QueryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Add(PostgresOverride, StatusName+"="+"500")
	path, status, ok = QueryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Add(PostgresOverride, PathName2+"="+"file://test/data")
	h.Add(PostgresOverride, StatusName+"="+"419")
	path, status, ok = QueryValues(h)
	fmt.Printf("test: QueryValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	//Output:
	//test: QueryValues() -> [path:] [status:0] [ok:false]
	//test: QueryValues() -> [path:] [status:0] [ok:false]
	//test: QueryValues() -> [path:] [status:0] [ok:false]
	//test: QueryValues() -> [path:file://test/data] [status:200] [ok:true]
	//test: QueryValues() -> [path:] [status:500] [ok:true]
	//test: QueryValues() -> [path:file://test/data] [status:419] [ok:true]

}
