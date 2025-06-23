package query

import (
	"bytes"
	"github.com/jackc/pgx/v5"
)

func marshal(rows pgx.Rows) (bytes.Buffer, error) {
	buf := bytes.Buffer{}
	if rows == nil {
		return buf, nil
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
		//err1 := fn(columnNames, values)
		//if err1 != nil {
		//	return buf,err1
		//}
	}
	return buf, nil
}
