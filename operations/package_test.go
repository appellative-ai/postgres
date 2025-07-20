package operations

import (
	"fmt"
	"github.com/appellative-ai/core/messaging"
	"github.com/appellative-ai/postgres/private"
	"time"
)

func ExampleConfigLogging() {
	private.FuncLogging(func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration) {
		fmt.Printf("test: ConfigLogging() -> %v\n", "logging direct invoke")
	})
	private.FuncLogging2(func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration) {
		fmt.Printf("test: ConfigLogging() -> %v\n", "logging direct invoke 2")
	})

	m := messaging.NewConfigMessage(func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration) {
		fmt.Printf("test: ConfigLogging() -> %v\n", "logging message invoke")
	})
	private.MessageLogging(m)

	private.MessageLogging2(m)

	var fn private.LogFunc2
	fn = func(start time.Time, duration time.Duration, req any, resp any, timeout time.Duration) {
		fmt.Printf("test: ConfigLogging() -> %v\n", "logging message invoke func")
	}
	m = messaging.NewConfigMessage(fn)
	private.MessageLogging2(m)

	//Output:
	//test: ConfigLogging() -> logging direct invoke
	//test: ConfigLogging() -> logging direct invoke 2
	//test: ConfigLogging() -> logging message invoke
	//test: MessageLogging2() -> logging message invoke failure
	//test: ConfigLogging() -> logging message invoke func
	
}
