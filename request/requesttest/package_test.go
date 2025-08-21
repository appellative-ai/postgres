package requesttest

import (
	"fmt"
	"github.com/appellative-ai/common/core"
)

func ExampleNewRequester() {
	m := core.NewSyncMap[string, any]()
	r := NewRequester(m)

	fmt.Printf("test: NewRequester() -> %v\n", r != nil)

	//Output:
	//test: NewRequester() -> true

}
