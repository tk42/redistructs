package redistructs

import (
	"context"

	"github.com/tk42/redistructs/types"
)

// RediStructs is an interface for operations for objects
type RediStructs interface {
	Get(ctx context.Context, dests ...types.RediStruct) error
	Put(ctx context.Context, src interface{}) error
	Delete(ctx context.Context, src interface{}) error
	// DeleteAll(ctx context.Context, mods ...rq.Modifier) error
	// Count(ctx context.Context, mods ...rq.Modifier) (int, error)
	// List(ctx context.Context, dest types.RediStruct, mods ...rq.Modifier) error

	// Map() map[string]types.RediStruct
	// Values() []types.RediStruct
	// Names() []string
	// IsZero() bool
	// HasZero() bool
}

// New creates a RediStructs instance
func New(pool Pool, config types.Config, model types.RediStruct) RediStructs {
	return NewRedigoStructs(pool, config, model)
}
