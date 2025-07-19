package private

import (
	"errors"
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
