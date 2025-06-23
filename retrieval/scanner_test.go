package retrieval

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/behavioral-ai/core/iox"
	json2 "github.com/behavioral-ai/core/json"
	"time"
)

const (
	EntriesPath = "file://[cwd]/querytest/entry.json"

	StartTimeName = "start_time"
	DurationName  = "duration_ms"
	TrafficName   = "traffic"
	CreatedTSName = "created_ts"

	RegionName  = "region"
	ZoneName    = "zone"
	SubZoneName = "sub_zone"
	HostName    = "host"

	MethodName     = "method"
	UrlName        = "url"
	PathName       = "path"
	StatusCodeName = "status_code"
	RouteName      = "route"
)

type testRows struct {
	index  int
	rows   []testEntry
	result []testEntry
}

func newTestRows(entries []testEntry) *testRows {
	e := new(testRows)
	if len(entries) > 0 {
		e.rows = append(e.rows, entries...)
	}
	e.index = -1
	return e
}

func (t *testRows) Close()     {}
func (t *testRows) Err() error { return nil }
func (t *testRows) Next() bool {
	if len(t.rows) == 0 || (t.index+1) >= len(t.rows) {
		return false
	}
	t.index++
	return true
}

func (t *testRows) Values() (result []any, err error) {
	result = append(result, t.rows[t.index].StartTime)
	result = append(result, t.rows[t.index].Duration)  //int64     `json:"duration"`
	result = append(result, t.rows[t.index].Traffic)   //string    `json:"traffic"`
	result = append(result, t.rows[t.index].CreatedTS) //time.Time `json:"created-ts"`

	result = append(result, t.rows[t.index].Region)  //string `json:"region"`
	result = append(result, t.rows[t.index].Zone)    //string `json:"zone"`
	result = append(result, t.rows[t.index].SubZone) //string `json:"sub-zone"`
	result = append(result, t.rows[t.index].Host)    //string `json:"host"`

	result = append(result, t.rows[t.index].Method)     //string `json:"method"`
	result = append(result, t.rows[t.index].Url)        //string `json:"url"`
	result = append(result, t.rows[t.index].Path)       //string `json:"path"`
	result = append(result, t.rows[t.index].StatusCode) //int32  `json:"status-code"`
	result = append(result, t.rows[t.index].Route)      //string `json:"route"`
	return
}

func (t *testRows) scan(columnNames []string, values []any) error {
	entry := testEntry{}
	for i, name := range columnNames {
		switch name {
		case StartTimeName:
			entry.StartTime = values[i].(time.Time)
		case DurationName:
			entry.Duration = values[i].(int64)
		case TrafficName:
			entry.Traffic = values[i].(string)
		case CreatedTSName:
			entry.CreatedTS = values[i].(time.Time)

		case RegionName:
			entry.Region = values[i].(string)
		case ZoneName:
			entry.Zone = values[i].(string)
		case SubZoneName:
			entry.SubZone = values[i].(string)
		case HostName:
			entry.Host = values[i].(string)

		case MethodName:
			entry.Method = values[i].(string)
		case UrlName:
			entry.Url = values[i].(string)
		case PathName:
			entry.Path = values[i].(string)
		case StatusCodeName:
			entry.StatusCode = values[i].(int32)
		case RouteName:
			entry.Route = values[i].(string)
		default:
			return errors.New(fmt.Sprintf("invalid field name: %v", name))
		}
	}
	t.result = append(t.result, entry)
	return nil
}

var (
	entries = []testEntry{
		{time.Now().UTC(), 100, "egress", time.Now().UTC(), "us-west", "oregon", "dc1", "www.test-host.com", "GET", "https://www.google.com/search?q-golang", "/search", 200, "google-search"},
		{time.Now().UTC(), 100, "egress", time.Now().UTC(), "us-central", "iowa", "dc1", "localhost:8081", "GET", "http://localhost:8081/advanced-go/search:google?q-golang", "/search", 200, "search"},
	}
	columnNames = []string{
		StartTimeName, DurationName, TrafficName, CreatedTSName,
		RegionName, ZoneName, SubZoneName, HostName,
		MethodName, UrlName, PathName, StatusCodeName, RouteName,
	}
)

// testEntry - timeseries access log struct
type testEntry struct {
	StartTime time.Time `json:"start-time"`
	Duration  int64     `json:"duration"`
	Traffic   string    `json:"traffic"`
	CreatedTS time.Time `json:"created-ts"`

	Region  string `json:"region"`
	Zone    string `json:"zone"`
	SubZone string `json:"sub-zone"`
	Host    string `json:"host"`

	Method     string `json:"method"`
	Url        string `json:"url"`
	Path       string `json:"path"`
	StatusCode int32  `json:"status-code"`
	Route      string `json:"route"`
}

func ExampleScanner() {
	rows := newTestRows(entries)
	err := Scanner(rows.scan, columnNames, rows)
	fmt.Printf("test: Scanner() -> [%v] [count:%v] [err:%v]\n", nil, len(rows.result), err)

	//Output:
	//test: Scanner() -> [<nil>] [count:2] [err:<nil>]

}

func _ExampleMarshal() {
	buf, err := json.Marshal(entries)
	fmt.Printf("test: json.Marshal() -> [buf:%v] [err:%v]\n", string(buf), err)

	buf, err = iox.ReadFile(EntriesPath)
	fmt.Printf("test: iox.ReadFile(\"%v\") -> [err:%v]\n", EntriesPath, err)

	list2, err2 := json2.New[[]testEntry](EntriesPath, nil)
	fmt.Printf("test: json2.New() -> [len:%v] [err:%v]\n", len(list2), err2)

	//Output:
	//test: json.Marshal() -> [buf:[{"start-time":"2025-06-22T13:58:47.845875Z","duration":100,"traffic":"egress","created-ts":"2025-06-22T13:58:47.845875Z","region":"us-west","zone":"oregon","sub-zone":"dc1","host":"www.test-host.com","method":"GET","url":"https://www.google.com/search?q-golang","path":"/search","status-code":200,"route":"google-search"},{"start-time":"2025-06-22T13:58:47.845875Z","duration":100,"traffic":"egress","created-ts":"2025-06-22T13:58:47.845875Z","region":"us-central","zone":"iowa","sub-zone":"dc1","host":"localhost:8081","method":"GET","url":"http://localhost:8081/advanced-go/search:google?q-golang","path":"/search","status-code":200,"route":"search"}]] [err:<nil>]
	//test: iox.ReadFile("file://[cwd]/querytest/entry.json") -> [err:<nil>]
	//test: json2.New() -> [len:2] [err:<nil>]

}

/*
	err := Scanner(func(columnNames []string, values []any) error {
		entry, err := scanEntry(columnNames, values)
		if err != nil {
			return err
		}
		result = append(result, entry)
		return nil
	}, columnNames, rows)

*/
