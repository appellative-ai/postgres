package commontest

import (
	"fmt"
	"github.com/behavioral-ai/postgres/common"
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
