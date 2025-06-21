package query

import "fmt"

func ExampleNewAgent() {
	a := newAgent()

	fmt.Printf("test: newAgent() -> [%v]\n", a)

	//Output:
	//test: newAgent() -> [sql:postgres:agent/query]

}
