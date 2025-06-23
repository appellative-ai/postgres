package retrieval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

func Marshaler(columnNames []string, rows Rows) (bytes.Buffer, error) {
	buf := bytes.Buffer{}
	if len(columnNames) == 0 {
		return buf, errors.New("column names list is empty")
	}
	if rows == nil {
		return buf, errors.New("rows list is empty")
	}
	var err error
	var values []any

	defer rows.Close()
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return buf, err
		}
		values, err = rows.Values()
		if err != nil || len(values) == 0 {
			return buf, err
		}
		// TODO: marshal values
		bytes, err1 := json.Marshal(values)
		if err1 != nil {
			return buf, err1
		}
		fmt.Printf("test: json.Marshal() -> [%v]\n", string(bytes))
	}
	return buf, nil
}
