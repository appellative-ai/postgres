package operations

import (
	"github.com/appellative-ai/core/messaging"
	"time"
)

const (
	userConfigKey = "user"
	pswdConfigKey = "pswd"
	uriConfigKey  = "uri"
	pingRouteName = "postgres-ping"
)

func ConfigClient(cfg map[string]string) error {
	return agent.clientStartup(cfg)
}

func ConfigLogging(log func(start time.Time, duration time.Duration, route string, req any, resp any, timeout time.Duration)) {
	agent.configLogging(log)
}

func Startup() {
	agent.Message(messaging.StartupMessage)

}

func Shutdown() {
	agent.Message(messaging.ShutdownMessage)
	clientShutdown()
}
