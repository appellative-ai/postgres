package exec

import (
	"errors"
	"fmt"
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

func ExampleInsert() {
	var h http.Header

	tag, status := InsertT[Entry](nil, h, "timeseries", "", nil, nil)
	fmt.Printf("test: Insert() -> [tag:%v] [status:%v]\n", tag, status)

	h = make(http.Header)
	tag, status = InsertT[Entry](nil, h, "timeseries", "", nil, nil)
	fmt.Printf("test: Insert() -> [tag:%v] [status:%v]\n", tag, status)

	h.Add(common.PostgresOverride, "count=321")
	tag, status = InsertT[Entry](nil, h, "timeseries", "", nil, nil)
	fmt.Printf("test: Insert() -> [tag:%v] [status:%v]\n", tag, status)

	h.Del(common.PostgresOverride)
	h.Add(common.PostgresOverride, "status=418")
	tag, status = InsertT[Entry](nil, h, "timeseries", "", nil, nil)
	fmt.Printf("test: Insert() -> [tag:%v] [status:%v]\n", tag, status)

	h.Add(common.PostgresOverride, "count=814")
	tag, status = InsertT[Entry](nil, h, "timeseries", "", nil, nil)
	fmt.Printf("test: Insert() -> [tag:%v] [status:%v]\n", tag, status)

	//Output:
	//test: Insert() -> [tag:{ 0 false false false false}] [status:Invalid Argument [DbClient is nil]]
	//test: Insert() -> [tag:{ 0 false false false false}] [status:Invalid Argument [DbClient is nil]]
	//test: Insert() -> [tag:{ 321 true false false false}] [status:OK]
	//test: Insert() -> [tag:{ 0 true false false false}] [status:I'm A Teapot]
	//test: Insert() -> [tag:{ 814 true false false false}] [status:I'm A Teapot]

}
