package common

import (
	"context"
	"fmt"
	"github.com/behavioral-ai/core/messaging"
	"net/http"
)

func (Entry) Get() ([]byte, *messaging.Status) {
	return nil, messaging.StatusNotFound()
}

func ExampleQueryT() {
	//e := Entry{}
	h := make(http.Header)

	rows, status := QueryT[Entry](context.Background(), h, "timeseries", "")

	fmt.Printf("test: Query() -> [rows:%v] [status:%v]\n", rows, status)

	//Output:
	//fail
}
