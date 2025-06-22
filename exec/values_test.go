package exec

import (
	"fmt"
	"github.com/behavioral-ai/postgres/common"
	"net/http"
)

func ExampleExeValues() {
	var h http.Header

	count, status, ok := execValues(h)
	fmt.Printf("test: execValues() -> [count:%v] [status:%v] [ok:%v]\n", count, status, ok)

	h = make(http.Header)
	h.Set(common.PostgresOverride, "test")
	count, status, ok = execValues(h)
	fmt.Printf("test: execValues() -> [count:%v] [status:%v] [ok:%v]\n", count, status, ok)

	h.Set(common.PostgresOverride, countName)
	count, status, ok = execValues(h)
	fmt.Printf("test: execValues() -> [count:%v] [status:%v] [ok:%v]\n", count, status, ok)

	h = make(http.Header)
	h.Add(common.PostgresOverride, countName+"="+"123")
	count, status, ok = execValues(h)
	fmt.Printf("test: execValues() -> [count:%v] [status:%v] [ok:%v]\n", count, status, ok)

	h = make(http.Header)
	h.Add(common.PostgresOverride, statusName+"="+"500")
	count, status, ok = execValues(h)
	fmt.Printf("test: execValues() -> [count:%v] [status:%v] [ok:%v]\n", count, status, ok)

	h = make(http.Header)
	h.Add(common.PostgresOverride, countName+"="+"456")
	h.Add(common.PostgresOverride, statusName+"="+"419")
	count, status, ok = execValues(h)
	fmt.Printf("test: execValues() -> [count:%v] [status:%v] [ok:%v]\n", count, status, ok)

	//Output:
	//test: execValues() -> [count:0] [status:0] [ok:false]
	//test: execValues() -> [count:0] [status:0] [ok:false]
	//test: execValues() -> [count:0] [status:0] [ok:false]
	//test: execValues() -> [count:123] [status:200] [ok:true]
	//test: execValues() -> [count:0] [status:500] [ok:true]
	//test: execValues() -> [count:456] [status:419] [ok:true]

}
