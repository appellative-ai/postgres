package query

import (
	"github.com/jackc/pgx/v5/pgconn"
)

// Scan - function for scanning rows
func scan(fn ScanFunc, columnNames []string, rows Rows) error {
	if rows == nil || len(columnNames) == 0 {
		return nil
	}
	var err error
	var values []any

	defer rows.Close()
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return err
		}
		values, err = rows.Values()
		if err != nil {
			return err
		}
		err1 := fn(columnNames, values)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func createColumnNames(fields []pgconn.FieldDescription) []string {
	var names []string
	for _, fld := range fields {
		names = append(names, fld.Name)
	}
	return names
}
