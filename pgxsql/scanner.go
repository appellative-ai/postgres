package pgxsql

import (
	"errors"
	"github.com/appellative-ai/core/jsonx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// Scanner - templated interface for scanning rows
type Scanner[T any] interface {
	Scan(columnNames []string, values []any) (T, error)
	Rows([]T) [][]any
}

// Unmarshal - templated function for JSON unmarshalling
func Unmarshal[T Scanner[T]](t any) ([]T, error) {
	if t == nil {
		return []T{}, errors.New("error: source is nil")
	}
	t2, err := jsonx.New[[]T](t, nil)
	//return json.New[[]T](t, nil)
	if err != nil {
		return t2, err //messaging.NewStatus(messaging.StatusJsonDecodeError, err)
	}
	return t2, nil //messaging.StatusOK()
}

// Rows - templated function for creating rows
func Rows[T Scanner[T]](entries []T) ([][]any, error) {
	if len(entries) == 0 {
		return nil, errors.New("not found") //messaging.StatusNotFound()
	}
	var t T
	return t.Rows(entries), nil //messaging.StatusOK()
}

// Scan - templated function for scanning rows
func Scan[T Scanner[T]](rows pgx.Rows) ([]T, error) {
	if rows == nil || rows.CommandTag().RowsAffected() == 0 {
		return nil, errors.New("not found") //messaging.StatusNotFound() //messaging.NewStatusError(messaging.StatusInvalidArgument, errors.New("invalid request: rows interface is nil"))
	}
	var s T
	var t []T
	var err error
	var values []any

	defer rows.Close()
	names := createColumnNames(rows.FieldDescriptions())
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return t, err
		}
		values, err = rows.Values()
		if err != nil {
			return t, err
		}
		val, err1 := s.Scan(names, values)
		if err1 != nil {
			return t, err
		}
		t = append(t, val)
		// Test this
		//rows.Close()
	}
	if len(t) == 0 {
		return t, errors.New("not found")
	}
	return t, nil
}

func createColumnNames(fields []pgconn.FieldDescription) []string {
	var names []string
	for _, fld := range fields {
		names = append(names, fld.Name)
	}
	return names
}
