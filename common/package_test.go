package common

import (
	"fmt"
	"net/http"
)

func ExampleLocationValues() {
	var h http.Header

	path, status, ok := LocationValues(h)
	fmt.Printf("test: LocationValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Set(LocationName, "test")
	path, status, ok = LocationValues(h)
	fmt.Printf("test: LocationValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h.Set(LocationName, PathName2)
	path, status, ok = LocationValues(h)
	fmt.Printf("test: LocationValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Add(LocationName, PathName2+"="+"file://test/data")
	path, status, ok = LocationValues(h)
	fmt.Printf("test: LocationValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Add(LocationName, StatusName+"="+"500")
	path, status, ok = LocationValues(h)
	fmt.Printf("test: LocationValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	h = make(http.Header)
	h.Add(LocationName, PathName2+"="+"file://test/data")
	h.Add(LocationName, StatusName+"="+"419")
	path, status, ok = LocationValues(h)
	fmt.Printf("test: LocationValues() -> [path:%v] [status:%v] [ok:%v]\n", path, status, ok)

	//Output:
	//test: LocationValues() -> [path:] [status:0] [ok:false]
	//test: LocationValues() -> [path:] [status:0] [ok:false]
	//test: LocationValues() -> [path:] [status:0] [ok:false]
	//test: LocationValues() -> [path:file://test/data] [status:200] [ok:true]
	//test: LocationValues() -> [path:] [status:500] [ok:true]
	//test: LocationValues() -> [path:file://test/data] [status:419] [ok:true]

}
