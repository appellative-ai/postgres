package operations

import (
	"fmt"
	"github.com/appellative-ai/postgres/request"
	"github.com/appellative-ai/postgres/retrieval"
)

func ExampleConfigLogging() {
	/*
		private.FuncLogging(func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration) {
			fmt.Printf("test: ConfigLogging() -> %v\n", "logging direct invoke")
		})
		private.FuncLogging2(func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration) {
			fmt.Printf("test: ConfigLogging() -> %v\n", "logging direct invoke 2")
		})

		m := messaging.NewConfigMessage(func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration) {
			fmt.Printf("test: ConfigLogging() -> %v\n", "logging message invoke")
		})
		private.MessageLogging(m)

		private.MessageLogging2(m)

		var fn private.LogFunc2
		fn = func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration) {
			fmt.Printf("test: ConfigLogging() -> %v\n", "logging message invoke func")
		}
		m = messaging.NewConfigMessage(fn)
		private.MessageLogging2(m)

	*/

	fmt.Printf("test: ConfigLogging()\n")

	//Output:
	//test: ConfigLogging()

}

type addressT struct {
	City  string
	State string
	Zip   string
}

func ExampleConfigSourceOverride() {
	name := "test:customer:type/address1"

	ConfigSourceOverride()
	r, err := request.Requester.Execute(nil, name, "insert", &addressT{
		City:  "Frisco",
		State: "TX",
		Zip:   "75035",
	})
	fmt.Printf("test: Requester.Execute() %v [err:%v]\n", r, err)

	buf, err1 := retrieval.Retriever.Marshal(nil, name, "get", nil)
	fmt.Printf("test: Retriever.Marshal() %v [err:%v]\n", string(buf.Bytes()), err1)

	buf, err1 = retrieval.Retriever.Marshal(nil, "invalid", "get", nil)
	fmt.Printf("test: Retriever.Marshal() %v [err:%v]\n", buf, err1)

	//Output:
	//test: Requester.Execute() { 1 true false false false} [err:<nil>]
	//test: Retriever.Marshal() {"City":"Frisco","State":"TX","Zip":"75035"} [err:<nil>]
	//test: Retriever.Marshal() <nil> [err:<nil>]

}
