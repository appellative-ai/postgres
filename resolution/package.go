package resolution

import (
	"github.com/appellative-ai/core/std"
)

// Interface - database access
type Interface struct {
	Representation    func(name string) (std.Content, *std.Status)
	AddRepresentation func(name, author, contentType string, value any) *std.Status

	Context    func(name string) (std.Content, *std.Status)
	AddContext func(name, author, contentType string, value any) *std.Status
}

// Resolver -
var Resolver = func() *Interface {
	return &Interface{
		Representation: func(name string) (std.Content, *std.Status) {
			return std.Content{}, std.StatusNotFound
		},
		AddRepresentation: func(name, author, contentType string, value any) *std.Status {
			return std.StatusOK
		},
		Context: func(name string) (std.Content, *std.Status) {
			return std.Content{}, std.StatusNotFound
		},
		AddContext: func(name, author, contentType string, t any) *std.Status {
			return std.StatusOK
		},
	}
}()
