package retrieval

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/behavioral-ai/core/fmtx"
	"reflect"
	"time"
)

const (
	textFmt     = "\"%v\":\"%v\""
	nonTextFmt  = "\"%v\":%v"
	arrayStart  = "["
	arrayEnd    = "]"
	objectStart = "{"
	objectEnd   = "}"
	endOfLine   = ","
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
	count := 0

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
		if count == 0 {
			buf.WriteString(arrayStart)
		}
		if count > 0 {
			buf.WriteString(endOfLine)
		}
		buf.WriteString(objectStart)
		writeValues(&buf, columnNames, values)
		buf.WriteString(objectEnd)
		count++
	}
	buf.WriteString(arrayEnd)
	return buf, nil
}

func writeValues(buf *bytes.Buffer, columnNames []string, values []any) {
	for i, v := range values {
		if i > 0 {
			//fmt.Printf(",")
			buf.WriteString(endOfLine)
		}
		t := reflect.TypeOf(v)
		switch t.Kind() {
		case reflect.String:
			s := fmt.Sprintf(textFmt, columnName(i, columnNames, v), v)
			//fmt.Printf(s)
			buf.WriteString(s)
		default:
			s := ""
			if ts, ok := v.(time.Time); ok {
				s = fmt.Sprintf(textFmt, columnName(i, columnNames, v), fmtx.FmtRFC3339Millis(ts))
			} else {
				s = fmt.Sprintf(nonTextFmt, columnName(i, columnNames, v), v)
			}
			//fmt.Printf(s)
			buf.WriteString(s)
		}
	}
}

func columnName(i int, names []string, v any) string {
	if i >= len(names) {
		t := reflect.TypeOf(v)
		return fmt.Sprintf("anonymous-%v-%v", i-len(names)+1, t.Name())
	}
	return names[i]
}
