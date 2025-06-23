package query

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

type entryRows struct {
	index int
	rows  []Entry
}

func newEntryRows(entries []Entry) *entryRows {
	e := new(entryRows)
	e.rows = append(e.rows, entries...)
	e.index = -1
	return e
}

func (e *entryRows) Close()     {}
func (e *entryRows) Err() error { return nil }
func (e *entryRows) Next() bool {
	if len(e.rows) == 0 || (e.index+1) >= len(e.rows) {
		return false
	}
	e.index++
	return true
}

func (e *entryRows) Values() (result []any, err error) {
	result = append(result, e.rows[e.index].StartTime)
	result = append(result, e.rows[e.index].Duration)  //int64     `json:"duration"`
	result = append(result, e.rows[e.index].Traffic)   //string    `json:"traffic"`
	result = append(result, e.rows[e.index].CreatedTS) //time.Time `json:"created-ts"`

	result = append(result, e.rows[e.index].Region)  //string `json:"region"`
	result = append(result, e.rows[e.index].Zone)    //string `json:"zone"`
	result = append(result, e.rows[e.index].SubZone) //string `json:"sub-zone"`
	result = append(result, e.rows[e.index].Host)    //string `json:"host"`

	result = append(result, e.rows[e.index].Method)     //string `json:"method"`
	result = append(result, e.rows[e.index].Url)        //string `json:"url"`
	result = append(result, e.rows[e.index].Path)       //string `json:"path"`
	result = append(result, e.rows[e.index].StatusCode) //int32  `json:"status-code"`
	result = append(result, e.rows[e.index].Route)      //string `json:"route"`
	return
}

func scanEntry(columnNames []string, values []any) (Entry, error) {
	entry := Entry{}
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
			return Entry{}, errors.New(fmt.Sprintf("invalid field name: %v", name))
		}
	}
	return entry, nil
}

var (
	list = []Entry{
		{time.Now().UTC(), 100, "egress", time.Now().UTC(), "us-west", "oregon", "dc1", "www.test-host.com", "GET", "https://www.google.com/search?q-golang", "/search", 200, "google-search"},
		{time.Now().UTC(), 100, "egress", time.Now().UTC(), "us-central", "iowa", "dc1", "localhost:8081", "GET", "http://localhost:8081/advanced-go/search:google?q-golang", "/search", 200, "search"},
	}
	columnNames = []string{
		StartTimeName, DurationName, TrafficName, CreatedTSName,
		RegionName, ZoneName, SubZoneName, HostName,
		MethodName, UrlName, PathName, StatusCodeName, RouteName,
	}
)

// Entry - timeseries access log struct
type Entry struct {
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

func ExampleScan() {
	var result []Entry
	rows := newEntryRows(list)
	err := scan(func(columnNames []string, values []any) error {
		entry, err := scanEntry(columnNames, values)
		if err != nil {
			return err
		}
		result = append(result, entry)
		return nil
	}, columnNames, rows)

	fmt.Printf("test: scan() -> [%v] [count:%v] [err:%v]\n", result != nil, len(result), err)

	//Output:
	//test: scan() -> [true] [count:2] [err:<nil>]

}

func _ExampleMarshal() {
	buf, err := json.Marshal(list)
	fmt.Printf("test: json.Marshal() -> [buf:%v] [err:%v]\n", string(buf), err)

	buf, err = iox.ReadFile(EntriesPath)
	fmt.Printf("test: iox.ReadFile(\"%v\") -> [err:%v]\n", EntriesPath, err)

	list2, err2 := json2.New[[]Entry](EntriesPath, nil)
	fmt.Printf("test: json2.New() -> [len:%v] [err:%v]\n", len(list2), err2)

	//Output:
	//test: json.Marshal() -> [buf:[{"start-time":"2025-06-22T13:58:47.845875Z","duration":100,"traffic":"egress","created-ts":"2025-06-22T13:58:47.845875Z","region":"us-west","zone":"oregon","sub-zone":"dc1","host":"www.test-host.com","method":"GET","url":"https://www.google.com/search?q-golang","path":"/search","status-code":200,"route":"google-search"},{"start-time":"2025-06-22T13:58:47.845875Z","duration":100,"traffic":"egress","created-ts":"2025-06-22T13:58:47.845875Z","region":"us-central","zone":"iowa","sub-zone":"dc1","host":"localhost:8081","method":"GET","url":"http://localhost:8081/advanced-go/search:google?q-golang","path":"/search","status-code":200,"route":"search"}]] [err:<nil>]
	//test: iox.ReadFile("file://[cwd]/querytest/entry.json") -> [err:<nil>]
	//test: json2.New() -> [len:2] [err:<nil>]

}
