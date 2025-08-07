package retrievaltest

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/appellative-ai/core/std"
	"github.com/appellative-ai/postgres/retrieval"
)

type retrievalT struct {
	cache *std.MapT[string, any]
}

func NewRetriever(m *std.MapT[string, any]) retrieval.Interface {
	r := new(retrievalT)
	r.cache = m
	return retrieval.Interface{Marshal: r.Marshal, Scan: r.Scan}
}

func (r *retrievalT) Marshal(ctx context.Context, name, sql string, args ...any) (bytes.Buffer, error) {
	t, ok := r.cache.Load(name)
	if !ok {
		return bytes.Buffer{}, nil
	}
	buf, err := json.Marshal(t)
	if err != nil {
		return bytes.Buffer{}, err
	}
	return *bytes.NewBuffer(buf), err
}

func (r *retrievalT) Scan(ctx context.Context, fn retrieval.ScanFunc, name, sql string, args ...any) error {
	//rows, err := agent.retrieve(ctx, name, sql, args)
	return nil //Scanner(fn, createColumnNames(rows.FieldDescriptions()), rows)
}
