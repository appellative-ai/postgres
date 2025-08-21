package requesttest

import (
	"fmt"
	"github.com/appellative-ai/core/std"
)

func ExampleNewRequester() {
	m := std.NewSyncMap[string, any]()
	r := NewRequester(m)

	fmt.Printf("test: NewRequester() -> %v\n", r != nil)

	//Output:
	//test: NewRequester() -> true

}
