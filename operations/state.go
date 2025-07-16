package operations

import (
	"github.com/appellative-ai/postgres/private"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type operationsT struct {
	running bool
	timeout time.Duration
	//Until    time.Duration
	log      private.LogFunc
	dbClient *pgxpool.Pool
}

// TODO: need to resolve all of the links in a collective and retrieval the registry for the
//
//	host names for the collective
//
//	Need a default domain for metadata/links -> root??, import??
/*
func initialize(m map[string]string) (ops *operationsT) {
	ops = new(operationsT)
	ops.registryHost1 = m[RegistryHost1Key]
	ops.registryHost2 = m[RegistryHost2Key]
	//ops.collective = cfg[CollectiveKey]
	//ops.domain = cfg[DomainKey]
	//var err error
	//if ops.origin, err = messaging.NewOriginFromMessage(msg); err != nil {
	//	// TODO: reply with error
	//	return
	//}
	//ops.serviceName = ops.collective + ":" + ops.origin.Name(ops.collective, ops.domain)
	return
}


*/
