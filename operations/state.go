package operations

type operationsT struct {
	running       bool
	dbClient pgxpool.

}

// TODO: need to resolve all of the links in a collective and query the registry for the
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