package diagnostic

import (
	"github.com/appellative-ai/core/messaging"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

func (a *agentT) config(m *messaging.Message) {
	if m == nil || m.Name != messaging.ConfigEvent {
		return
	}
	if a.running.Load() {
		return
	}
	if messaging.UpdateContent[func(start time.Time, duration time.Duration, route string, req any, resp any, timeout time.Duration)](m, &a.logFunc) {
		return
	}
	messaging.UpdateContent[*pgxpool.Pool](m, &a.dbClient)
}
