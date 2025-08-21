package operations

import (
	"github.com/appellative-ai/common/core"
	"github.com/appellative-ai/common/messaging"
	"github.com/appellative-ai/postgres/request"
	"github.com/appellative-ai/postgres/request/requesttest"
	"github.com/appellative-ai/postgres/retrieval"
	"github.com/appellative-ai/postgres/retrieval/retrievaltest"
	"time"
)

const (
	userConfigKey = "user"
	pswdConfigKey = "pswd"
	uriConfigKey  = "uri"
	pingRouteName = "postgres-ping"
)

var (
	cache *core.MapT[string, any]
)

func ConfigClient(cfg map[string]string) error {
	return agent.clientStartup(cfg)
}

func ConfigLogging(log func(start time.Time, duration time.Duration, route string, req any, resp any, timeout time.Duration)) {
	agent.configLogging(log)
}

func ConfigSourceOverride() {
	if cache == nil {
		cache = core.NewSyncMap[string, any]()
		retrieval.Retriever = retrievaltest.NewRetriever(cache)
		request.Requester = requesttest.NewRequester(cache)
	}
}

func Startup() {
	agent.Message(messaging.StartupMessage)

}

func Shutdown() {
	agent.Message(messaging.ShutdownMessage)
	clientShutdown()
}
