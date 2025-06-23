package query

import "fmt"

func ExampleMarshaler() {
	buf, err := Marshaler(nil, nil)
	fmt.Printf("test: Marshaler() [bytes:%v] [%v]\n", buf.Len(), err)

	buf, err = Marshaler(columnNames, nil)
	fmt.Printf("test: Marshaler() [bytes:%v] [%v]\n", buf.Len(), err)

	buf, err = Marshaler(columnNames, newTestRows(entries))
	fmt.Printf("test: Marshaler() [bytes:%v] [%v]\n", buf.Len(), err)

	//Output:
	//test: Marshaler() [bytes:0] [column names list is empty]
	//test: Marshaler() [bytes:0] [rows list is empty]
	//test: Marshaler() [bytes:0] [<nil>]

}
