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

	Map(ctx context.Context) (map[string]types.RediStruct, error)
	Values(ctx context.Context) ([]types.RediStruct, error)
	Names(ctx context.Context) ([]string, error)
	// IsZero(ctx context.Context) (bool, error)
	// HasZero(ctx context.Context) (bool, error)

	clone() interface{}
}

// New creates a RediStructs instance
func New(pool Pool, config types.Config, model types.RediStruct) RediStructs {
	return NewRedigoStructs(pool, config, model)
}
