package common

import (
	"encoding/json"
	"fmt"

	"time"
)

type accessLogV2 struct {
	Duration string
}

var list = []Entry{
	{time.Now().UTC(), 100, "egress", time.Now().UTC(), "us-west", "oregon", "dc1", "www.test-host.com", "123456", "req-id", "relate-to", "HTTP/1.1", "GET", "www.google.com", "", "https://www.google.com/search?q-golang", "/search", 200, "gzip", 12345, "google-search", "primary", 500, 98.5, 10, "RL"},
	{time.Now().UTC(), 100, "egress", time.Now().UTC(), "us-west", "oregon", "dc1", "localhost:8081", "123456", "req-id", "relate-to", "HTTP/1.1", "GET", "github/advanced-go/search", "", "http://localhost:8081/advanced-go/search:google?q-golang", "/search", 200, "gzip", 12345, "search", "primary", 500, 100, 10, "TO"},
}

func ExampleEntry() {
	buf, err := json.Marshal(list)
	if err != nil {
		fmt.Printf("test: Entry{} -> [err:%v]\n", err)
	} else {
		fmt.Printf("test: Entry{} -> %v\n", string(buf))
	}
	entry := Entry{}

	if v, ok := any(entry).(Variant); ok {
		if v != nil {
		}
		fmt.Printf("test: Query() -> [variant:%v]\n", ok)
	}
	//Output:
	//fail

}

func ExampleScanColumnsTemplate() {
	//log := scanColumnsTemplate[AccessLog](nil)

	//fmt.Printf("test: scanColumnsTemplate[AccessLog](nil) -> %v\n", log)

	//Output:
	//fail
}

func _ExampleScannerInterface_V1() {

	//log, status := scanRowsTemplateV1[AccessLog, AccessLog](nil)
	//fmt.Printf("test: scanRowsTemplateV1() -> [status:%v] [elem:%v] [log:%v] \n", status, reflect.TypeOf(log).Elem(), log[0].DurationString)

	//Output:
	//test: scanRowsTemplateV1() -> [status:OK] [elem:timeseries.AccessLog] [log:SCAN() TEST DURATION STRING]

}

func _ExampleScannerInterface() {
	//log, status := scanRowsTemplate[accessLogV2](nil)

	//log, status := scanRowsTemplate[AccessLog](nil)
	//fmt.Printf("test: scanRowsTemplate() -> [status:%v] [elem:%v] [log:%v] \n", status, reflect.TypeOf(log).Elem(), log[0].DurationString)

	//Output:
	//test: scanRowsTemplateV1() -> [status:OK] [elem:timeseries.AccessLog] [log:SCAN() TEST DURATION STRING]

}
