package pgxsql

import (
	"errors"
	"fmt"
	"github.com/appellative-ai/postgres/pgxdml"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

const (
	status504Q = "file://[cwd]/test/status-504.json"
)

type TestConditions struct {
	Time        time.Time
	Location    string
	Temperature float64
}

type rowsT struct {
}

func (r *rowsT) CommandTag() pgconn.CommandTag {
	//TODO implement me
	return pgconn.CommandTag{}
}

func (r *rowsT) FieldDescriptions() []pgconn.FieldDescription {
	//TODO implement me
	return nil
}

func (r *rowsT) Conn() *pgx.Conn {
	//TODO implement me
	return nil
}

func (r *rowsT) Close()     {}
func (r *rowsT) Err() error { return nil }

//	func (r *rowsT) CommandTag() CommandTag {
//		return pgconn.CommandTag{}//RowsAffected: 1, Insert: false, Update: false, Delete: false, Select: true}
//	}
//
// func (r *rowsT) FieldDescriptions() []FieldDescription { return nil }
func (r *rowsT) Next() bool             { return false }
func (r *rowsT) Scan(dest ...any) error { return nil }
func (r *rowsT) Values() ([]any, error) { return nil, nil }
func (r *rowsT) RawValues() [][]byte    { return nil }

const (
	queryErrorSql = "select * from test"
	queryRowsSql  = "select * from table"

	queryConditions      = "select * from conditions"
	queryConditionsWhere = "select * from conditions where $1 order by temperature desc"
	queryConditionsError = "select test,test2 from conditions"
	queryErrorRsc        = "error"
	queryRowsRsc         = "rows"
)

func ExampleQuery_TestError() {
	result, status := query(nil, newQueryRequest(queryErrorRsc, queryErrorSql, nil))
	fmt.Printf("test: retrieval(nil,%v) -> [rows:%v] [status:%v]\n", queryErrorSql, result, status)

	//Output:
	//test: retrieval(nil,select * from test) -> [rows:<nil>] [status:Invalid Argument [error on PostgreSQL database retrieval call: dbClient is nil]]

}

func ExampleQuery_StatusTimeout() {
	//setOverrideLookup([]string{"", status504Q})
	//lookup.SetOverride(status504Q)
	rows, status := query(nil, newQueryRequest(queryRowsRsc, queryRowsSql, nil))
	fmt.Printf("test: retrieval(nil,%v) -> [rows:%v] [status:%v]\n", queryRowsSql, rows, status)

	//Output:
	//test: retrieval(nil,select * from table) -> [rows:<nil>] [status:Timeout [error 1]]

}

func ExampleQuery_Proxy() {
	// Need to clear per test override
	//lookup.SetOverride(io2.StatusOKUri) //setOverrideLookup([]string{"", ""})
	req := newQueryRequest(queryRowsRsc, queryRowsSql, nil)
	rows, status := query(nil, req)
	fmt.Printf("test: retrieval(ctx,%v) -> [rows:%v] [status:%v]\n", queryRowsSql, rows, status)

	//Output:
	//test: retrieval(ctx,select * from table) -> [rows:<nil>] [status:OK]

}

func ExampleQuery_Conditions_Error() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer clientShutdown()
		req := newQueryRequest(queryRowsRsc, queryConditionsError, nil)
		rows, status := query(nil, req)
		if status != nil {
			fmt.Printf("test: retrieval(nil,%v) -> [status:%v]\n", queryConditionsError, status)
		} else {
			fmt.Printf("test: retrieval(nil,%v) -> [status:%v] [cmd:%v]\n", queryConditions, status, rows.CommandTag())
			conditions, status1 := processResults(rows, "")
			fmt.Printf("test: processResults(results) -> [status:%v] [rows:%v]\n", status1, conditions)
		}
	}

	//Output:
	//[[] github.com/gotemplates/postgresql/pgxsql/retrieval [serverity:ERROR, code:42703, message:column "test" does not exist, position:8, SQLState:42703]]
	//test: retrieval(nil,select test,test2 from conditions) -> [status:Internal]

}

func ExampleQuery_Conditions() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer clientShutdown()
		req := newQueryRequest(queryRowsRsc, queryConditions, nil)
		results, status := query(nil, req)
		if status != nil {
			fmt.Printf("test: retrieval(nil,%v) -> [status:%v]\n", queryConditions, status)
		} else {
			fmt.Printf("test: retrieval(nil,%v) -> [status:%v] [cmd:%v]\n", queryConditions, status, results.CommandTag())
			conditions, status1 := processResults(results, "")
			fmt.Printf("test: processResults(results) -> [status:%v] [rows:%v]\n", status1, conditions)
		}
	}

	//Output:
	//test: retrieval(nil,select * from conditions) -> [status:OK] [cmd:{ 0 false false false false}]
	//test: processResults(results) -> [status:OK] [rows:[{2023-01-26 12:09:12.426535 -0600 CST office 70} {2023-01-26 12:09:12.426535 -0600 CST basement 66.5} {2023-01-26 12:09:12.426535 -0600 CST garage 45.1234}]]

}

func ExampleQuery_Conditions_Where() {
	err := testStartup()
	if err != nil {
		fmt.Printf("test: testStartup() -> [error:%v]\n", err)
	} else {
		defer clientShutdown()

		where := []pgxdml.Attr{{"location", "garage"}}
		req := newQueryRequest(queryRowsRsc, queryConditionsWhere, where)
		results, status := query(nil, req)
		if status != nil {
			fmt.Printf("test: retrieval(nil,%v) -> [status:%v]\n", queryConditionsWhere, status)
		} else {
			fmt.Printf("test: retrieval(nil,%v) -> [status:%v] [cmd:%v]\n", queryConditions, status, results.CommandTag())
			conditions, status1 := processResults(results, "")
			fmt.Printf("test: processResults(results) -> [status:%v] [rows:%v]\n", status1, conditions)
		}
	}

	//Output:
	//test: retrieval(nil,select * from conditions) -> [status:OK] [cmd:{ 0 false false false false}]
	//test: processResults(results) -> [status:OK] [rows:[{2023-01-26 12:09:12.426535 -0600 CST garage 45.1234}]]

}

func processResults(results pgx.Rows, msg string) (conditions []TestConditions, status error) {
	conditions, status = scanRows(results)
	if status == nil && len(conditions) == 0 {
		return nil, nil
	}
	return conditions, status
}

func scanRows(rows pgx.Rows) ([]TestConditions, error) {
	if rows == nil {
		return nil, errors.New("invalid request: Rows interface is empty")
	}
	var err error
	var values []any
	var conditions []TestConditions
	for rows.Next() {
		err = rows.Err()
		if err != nil {
			return nil, err
		}
		values, err = rows.Values()
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, scanColumns(values))
	}
	return conditions, nil
}

func scanColumns(values []any) TestConditions {
	var ts = new(pgtype.Timestamp)
	ts.Scan(values[0])

	cond := TestConditions{
		Time:        ts.Time,
		Location:    values[1].(string),
		Temperature: values[2].(float64),
	}
	return cond
}
