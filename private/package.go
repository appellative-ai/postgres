package private

import (
	"errors"
	"fmt"
	"github.com/appellative-ai/core/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

const (
	ContentTypeConfiguration = "application/x-configuration"
	ThresholdRequest         = "x-threshold-request"
	ThresholdResponse        = "x-threshold-response"
	ThresholdTimeoutName     = "timeout"

	TrafficEgress = "egress"
)

type LogFunc func(traffic string, start time.Time, duration time.Duration, route string, req any, resp any)

// Configuration -
type Configuration struct {
	//Timeout  time.Duration
	//Until    time.Duration
	Log      LogFunc
	DbClient *pgxpool.Pool
}

func (c *Configuration) Update(cfg *Configuration) {

	if cfg.Log != nil {
		c.Log = cfg.Log
	}
	if cfg.DbClient != nil {
		c.DbClient = cfg.DbClient
	}
}

func NewConfiguration(timeout time.Duration) *Configuration {
	c := new(Configuration)
	return c
}

func NewConfigurationMessage(cfg *Configuration) *messaging.Message {
	m := messaging.NewMessage(messaging.ChannelControl, messaging.ConfigEvent)
	m.SetContent(ContentTypeConfiguration, cfg)
	return m
}

func ConfigurationContent(m *messaging.Message) (*Configuration, *messaging.Status) {
	if !messaging.ValidContent(m, messaging.ConfigEvent, ContentTypeConfiguration) {
		return nil, messaging.NewStatus(messaging.StatusInvalidContent, errors.New("invalid content"))
	}
	return messaging.New[*Configuration](m.Content)
}

type LogFunc2 func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration)

func FuncLogging2(log LogFunc2) {
	log(time.Now(), 0, nil, nil, 0)

}
func FuncLogging(log func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration)) {
	log(time.Now(), 0, nil, nil, 0)
}

func MessageLogging(m *messaging.Message) {
	var fn func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration)
	ok := messaging.UpdateContent[func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration)](&fn, m)
	if ok {
		fn(time.Now(), 0, nil, nil, 0)
	}
}

func MessageLogging2(m *messaging.Message) {
	var fn LogFunc2
	ok := messaging.UpdateContent[LogFunc2](&fn, m)
	if ok {
		fn(time.Now(), 0, nil, nil, 0)
	} else {
		fmt.Printf("test: MessageLogging2() -> %v\n", "logging message invoke failure")
	}
}
