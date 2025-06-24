package retrieval

import (
	"context"
	"errors"
	"fmt"
)

const (
	timeseriesResource = "timeseries"
)

func ExampleTestScan() {
	err := testScan(nil, nil, "", "")
	fmt.Printf("test: Scan() -> [%v]\n", err)

	err = testScan(nil, nil, "invalid-resource", "")
	fmt.Printf("test: Scan() -> [%v]\n", err)

	err = testScan(nil, nil, timeseriesResource, "")
	fmt.Printf("test: Scan() -> [%v]\n", err)

	rows := newTestRows(nil)
	err = testScan(nil, rows.scan, timeseriesResource, "")
	fmt.Printf("test: Scan() -> [count:%v] [%v]\n", len(rows.result), err)

	//Output:
	//test: Scan() -> [resource not supported : ]
	//test: Scan() -> [resource not supported : invalid-resource]
	//test: Scan() -> [scanner ScanFunc is nil]
	//test: Scan() -> [count:2] [<nil>]

}

// Scan -  process a SQL select statement using a Scan function
func testScan(ctx context.Context, fn ScanFunc, resource, sql string, args ...any) error {
	switch resource {
	case timeseriesResource:
		return Scanner(fn, columnNames, newTestRows(entries))
	}
	return errors.New(fmt.Sprintf("resource not supported : %v", resource))
}

func process(relation *Resolution) error {
	rows := newTestRows(nil)
	return relation.Scan(nil, rows.scan, timeseriesResource, "")
}
