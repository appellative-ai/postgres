package requesttest

import (
	"context"
	"errors"
	"fmt"
	"github.com/appellative-ai/core/std"
	"github.com/appellative-ai/postgres/request"
	"strings"
)

const (
	insertSql = "insert"
)

type requestT struct {
	cache *std.MapT[string, any]
}

func NewRequester(m *std.MapT[string, any]) request.Interface {
	r := new(requestT)
	r.cache = m
	return request.Interface{Execute: r.Execute}
}

func (r *requestT) Execute(ctx context.Context, name, sql string, args ...any) (request.Result, error) {
	if name == "" || len(args) == 0 {
		return request.Result{}, errors.New(fmt.Sprintf("name is empty [%v] or args are empty", name))
	}
	if strings.Contains(name, insertSql) {
		r.cache.Store(name, args[0])
		return request.Result{Insert: true, RowsAffected: 1}, nil
	}
	return request.Result{}, errors.New(fmt.Sprintf("invalid SQL [%v]", sql))
}
