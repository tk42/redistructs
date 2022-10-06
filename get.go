package redistructs

import (
	"context"
	"reflect"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/tk42/redistructs/types"
)

func (rs *RedigoStructs) Get(ctx context.Context, dests ...types.RediStruct) error {
	conn, err := rs.pool.GetContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to acquire a connection")
	}
	defer conn.Close()

	dbIdx, changed := rs.getDBIndex(rs.model)
	if changed {
		conn.Do("SELECT", dbIdx)
	}

	for _, dest := range dests {
		if reflect.TypeOf(rs.model) != reflect.TypeOf(dest) {
			return errors.New("failed to compare RediStruct")
		}

		switch dest.StoreType() {
		case types.Serialized:
			b, err := dest.Marshal()
			if err != nil {
				return errors.Wrap(err, "failed to marshal")
			}
			if len(b) == 0 {
				return errors.Errorf("failed to marshal due to empty %v", dest)
			}
			err = rs.scripts["HGETALLXP"].Send(conn, rs.name, dest.PrimaryKey())
			if err != nil {
				return errors.Wrapf(err, "failed to send HGETALLXP %s", rs.key)
			}
		default:
			panic("unsupported store type")
		}
	}

	err = conn.Flush()
	if err != nil {
		return errors.Wrap(err, "faild to flush commands")
	}

	for _, dest := range dests {
		v, err := redigo.Values(conn.Receive())
		if err != nil {
			return errors.Wrap(err, "faild to receive or cast redis command result")
		}

		switch dest.StoreType() {
		case types.Serialized:
			var key string
			for j, vv := range v {
				if j%2 == 0 {
					key = string(vv.([]byte))
					continue
				} else if key == dest.PrimaryKey() {
					if err := dest.Unmarshal(vv.([]byte)); err != nil {
						return errors.Wrap(err, "failed to unmarshal")
					}
					break
				}
			}
		default:
			panic("unsupported store type")
		}
	}

	return nil
}
