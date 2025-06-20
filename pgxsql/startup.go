package pgxsql

import (
	"github.com/behavioral-ai/core/messaging"
	"sync/atomic"
)

var (
	ready int64
	agent *messaging.Agent
)

func isReady() bool {
	return atomic.LoadInt64(&ready) != 0
}

func setReady() {
	atomic.StoreInt64(&ready, 1)
}

func resetReady() {
	atomic.StoreInt64(&ready, 0)
}

func init() {
	//a, err1 := host.RegisterControlAgent(PkgPath, messageHandler)
	//if err1 != nil {
	//	fmt.Printf("init(\"%v\") failure: [%v]\n", PkgPath, err1)
	//	}
	//	a.Run()
}

/*
func messageHandler(msg *messaging.Message) {
	switch msg.Event() {
	case messaging.StartupEvent:
		start := time.Now()
		// TODO
		//clientStartup(msg)
		messaging.SendReply(msg, messaging.NewStatusDuration(http.StatusOK, time.Since(start)))
	case messaging.ShutdownEvent:
		clientShutdown()
	case messaging.PingEvent:
		start := time.Now()
		messaging.Reply(msg, messaging.NewStatusDuration(http.StatusOK, time.Since(start)))
	}
}


*/
