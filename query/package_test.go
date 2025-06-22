package query

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/behavioral-ai/core/iox"
	json2 "github.com/behavioral-ai/core/json"
	"github.com/behavioral-ai/postgres/common"
	"net/http"
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

var list = []Entry{
	{time.Now().UTC(), 100, "egress", time.Now().UTC(), "us-west", "oregon", "dc1", "www.test-host.com", "GET", "https://www.google.com/search?q-golang", "/search", 200, "google-search"},
	{time.Now().UTC(), 100, "egress", time.Now().UTC(), "us-central", "iowa", "dc1", "localhost:8081", "GET", "http://localhost:8081/advanced-go/search:google?q-golang", "/search", 200, "search"},
}

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

func (Entry) Scan(columnNames []string, values []any) (log Entry, err error) {
	for i, name := range columnNames {
		switch name {
		case StartTimeName:
			log.StartTime = values[i].(time.Time)
		case DurationName:
			log.Duration = values[i].(int64)
		case TrafficName:
			log.Traffic = values[i].(string)
		case CreatedTSName:
			log.CreatedTS = values[i].(time.Time)

		case RegionName:
			log.Region = values[i].(string)
		case ZoneName:
			log.Zone = values[i].(string)
		case SubZoneName:
			log.SubZone = values[i].(string)
		case HostName:
			log.Host = values[i].(string)

		case MethodName:
			log.Method = values[i].(string)
		case UrlName:
			log.Url = values[i].(string)
		case PathName:
			log.Path = values[i].(string)
		case StatusCodeName:
			log.StatusCode = values[i].(int32)
		case RouteName:
			log.Route = values[i].(string)

		default:
			err = errors.New(fmt.Sprintf("invalid field name: %v", name))
			return
		}
	}
	return
}

func (a Entry) Values() []any {
	return []any{
		a.StartTime,
		a.Duration,
		a.Traffic,
		a.CreatedTS,

		a.Region,
		a.Zone,
		a.SubZone,
		a.Host,

		a.Method,
		a.Url,
		a.Path,
		a.StatusCode,
		a.Route,
	}
}

func (Entry) Rows(events []Entry) [][]any {
	var values [][]any

	for _, e := range events {
		values = append(values, e.Values())
	}
	return values
}

func _ExampleEntry() {
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

func ExampleQuery() {
	h := make(http.Header)

	r, status := QueryT[Entry](context.Background(), h, "timeseries", "")
	fmt.Printf("test: Query() -> [count:%v] [status:%v]\n", len(r), status)

	h.Add(common.PostgresOverride, "status=418")
	r, status = QueryT[Entry](context.Background(), h, "timeseries", "")
	fmt.Printf("test: Query() -> [count:%v] [status:%v]\n", len(r), status)

	h.Del(common.PostgresOverride)
	h.Add(common.PostgresOverride, "path="+EntriesPath)
	r, status = QueryT[Entry](context.Background(), h, "timeseries", "")
	fmt.Printf("test: Query() -> [count:%v] [rows:%v] [status:%v]\n", len(r), r != nil, status)

	//Output:
	//test: Query() -> [count:0] [status:Internal Error [DbClient is nil]]
	//test: Query() -> [count:0] [status:I'm A Teapot]
	//test: Query() -> [count:2] [rows:true] [status:OK]

}
