package commontest

import (
	"fmt"
	"github.com/appellative-ai/postgres/common"
)

func ExampleEntryVariant() {
	e := common.Entry{}
	if v, ok := any(e).(common.Variant); ok {
		if v != nil {
		}
		fmt.Printf("test: Query() -> [variant:%v]\n", ok)
	}

	//Output:
	//fail
}
