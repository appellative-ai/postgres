package operations

import (
	"fmt"
)

func ExampleClientStartup() {
	//rsc := startupResource{Uri: ""}
	err := clientStartup(nil)
	if err != nil {
		defer clientShutdown()
	}
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	err = clientStartup(make(map[string]string))
	fmt.Printf("test: ClientStartup() -> %v\n", err)

	//Output:
	//test: ClientStartup() -> error: strings map configuration is nil
	//test: ClientStartup() -> database URL is empty

}
