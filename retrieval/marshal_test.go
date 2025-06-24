package retrieval

import (
	bytes2 "bytes"
	"encoding/json"
	"fmt"
)

func _ExampleMarshaler_Error() {
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

func ExampleMarshaler() {
	var bytes bytes2.Buffer

	buf, err := Marshaler(columnNames, newTestRows(entries))
	json.Indent(&bytes, buf.Bytes(), "", "")
	fmt.Printf("test: Marshaler() %v [%v]\n", buf.String(), err)

	buf, err = Marshaler(columnNames[:len(columnNames)-2], newTestRows(entries))
	json.Indent(&bytes, buf.Bytes(), "", "")
	fmt.Printf("test: Marshaler() %v [%v]\n", bytes.String(), err)

	//Output:
	//test: Marshaler() [bytes:0] [column names list is empty]

}
