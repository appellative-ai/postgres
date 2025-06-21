package common

import (
	"errors"
	"github.com/behavioral-ai/core/json"
	"github.com/behavioral-ai/core/messaging"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Variant interface {
	Get() ([]byte, *messaging.Status)
}

// Scanner - templated interface for scanning rows
type Scanner[T any] interface {
	Scan(columnNames []string, values []any) (T, error)
	Rows([]T) [][]any
}

// Unmarshal - templated function for JSON unmarshalling
func Unmarshal[T Scanner[T]](t any) ([]T, *messaging.Status) {
	if t == nil {
		return []T{}, messaging.NewStatus(messaging.StatusInvalidArgument, errors.New("error: source is nil"))
	}
	t2, err := json.New[[]T](t, nil)
	if err != nil {
		return t2, messaging.NewStatus(messaging.StatusJsonDecodeError, err)
	}
	return t2, messaging.StatusOK()
}

// Rows - templated function for creating rows
func Rows[T Scanner[T]](entries []T) ([][]any, *messaging.Status) {
	if len(entries) == 0 {
		return nil, messaging.StatusNotFound()
	}
	var t T
	return t.Rows(entries), messaging.StatusOK()
}

// Scan - templated function for scanning rows
func Scan[T Scanner[T]](rows pgx.Rows) ([]T, *messaging.Status) {
	if rows == nil || rows.CommandTag().RowsAffected() == 0 {
		return nil, messaging.StatusNotFound() //messaging.NewStatusError(messaging.StatusInvalidArgument, errors.New("invalid request: rows interface is nil"))
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
			return t, messaging.NewStatus(messaging.StatusInvalidArgument, err)
		}
		values, err = rows.Values()
		if err != nil {
			return t, messaging.NewStatus(messaging.StatusInvalidArgument, err)
		}
		val, err1 := s.Scan(names, values)
		if err1 != nil {
			return t, messaging.NewStatus(messaging.StatusInvalidArgument, err1)
		}
		t = append(t, val)
		// Test this
		//rows.Close()
	}
	if len(t) == 0 {
		return t, messaging.StatusNotFound()
	}
	return t, messaging.StatusOK()
}

func createColumnNames(fields []pgconn.FieldDescription) []string {
	var names []string
	for _, fld := range fields {
		names = append(names, fld.Name)
	}
	return names
}
