package retrievaltest

import (
	"fmt"
	"github.com/appellative-ai/core/std"
)

func ExampleNewRetriever() {
	m := std.NewSyncMap[string, any]()
	r := NewRetriever(m)

	fmt.Printf("test: NewRetriever() -> %v\n", r != nil)

	//Output:
	//test: NewRetriever() -> true

}
