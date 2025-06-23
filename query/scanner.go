package query

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

// Scanner - function for scanning rows
func Scanner(fn ScanFunc, columnNames []string, rows Rows) error {
	if fn == nil {
		return errors.New("scanner ScanFunc is nil")
	}
	if len(columnNames) == 0 {
		return errors.New("list of column names is empty")
	}
	if rows == nil {
		return errors.New("rows list is empty")
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
