package retrievaltest

import (
	"fmt"
	"github.com/appellative-ai/common/core"
)

func ExampleNewRetriever() {
	m := core.NewSyncMap[string, any]()
	r := NewRetriever(m)

	fmt.Printf("test: NewRetriever() -> %v\n", r != nil)

	//Output:
	//test: NewRetriever() -> true

}
