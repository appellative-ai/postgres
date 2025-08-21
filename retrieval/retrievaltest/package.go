package retrievaltest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/appellative-ai/core/std"
	"github.com/appellative-ai/postgres/retrieval"
	"strings"
)

const (
	get   = "get"
	query = "query"
)

type retrievalT struct {
	cache *std.MapT[string, any]
}

func NewRetriever(m *std.MapT[string, any]) *retrieval.Interface {
	r := new(retrievalT)
	r.cache = m
	return &retrieval.Interface{Marshal: r.Marshal, Scan: r.Scan}
}

func (r *retrievalT) Marshal(ctx context.Context, name, sql string, args ...any) (*bytes.Buffer, error) {
	s := strings.ToLower(sql)

	if strings.Contains(s, get) {
		t, ok := r.cache.Load(name)
		if !ok {
			return nil, nil
		}
		buf, err := json.Marshal(t)
		if err != nil {
			return nil, err
		}
		return bytes.NewBuffer(buf), err
	}
	return nil, errors.New(fmt.Sprintf("invalid SQL: %v", sql))
}

func (r *retrievalT) Scan(ctx context.Context, fn retrieval.ScanFunc, name, sql string, args ...any) error {
	//rows, err := agent.retrieve(ctx, name, sql, args)
	return nil //Scanner(fn, createColumnNames(rows.FieldDescriptions()), rows)
}
