package resolution

import (
	"context"
	"errors"
	"github.com/appellative-ai/core/std"
)

// Interface - database access
type Interface struct {
	Representation    func(ctx context.Context, name string) (std.Content, error)
	AddRepresentation func(ctx context.Context, name, author, contentType string, value any) error

	Context    func(ctx context.Context, name string) (std.Content, error)
	AddContext func(ctx context.Context, name, author, contentType string, value any) error
}

// Resolver -
var Resolver = func() *Interface {
	return &Interface{
		Representation: func(ctx context.Context, name string) (std.Content, error) {
			return std.Content{}, errors.New("not implemented")
		},
		AddRepresentation: func(ctx context.Context, name, author, contentType string, value any) error {
			return errors.New("not implemented")
		},
		Context: func(ctx context.Context, name string) (std.Content, error) {
			return std.Content{}, errors.New("not implemented")
		},
		AddContext: func(ctx context.Context, name, author, contentType string, t any) error {
			return errors.New("not implemented")
		},
	}
}()
