package redistructs

import (
	"context"
	"reflect"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/tk42/redistructs/types"
)

func (rs RedigoStructs) Get(ctx context.Context, dests ...types.RediStruct) error {
	conn, err := rs.pool.GetContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to acquire a connection")
	}
	defer conn.Close()

	dbIdx, changed := rs.getDBIndex(rs.model)
	if changed {
		conn.Do("SELECT", dbIdx)
	}

	_, err = conn.Do("WATCH", rs.name)
	if err != nil {
		return errors.Wrapf(err, "failed to send WATCH %ss", rs.name)
	}

	err = conn.Send("MULTI")
	if err != nil {
		return errors.Wrap(err, "faild to send MULTI command")
	}

	for _, dest := range dests {
		if reflect.ValueOf(rs.model).Elem() != reflect.ValueOf(dest).Elem() {
			return errors.New("failed to compare RediStruct")
		}

		switch dest.StoreType() {
		case types.Serialized:
			if len(dest.Serialized()) == 0 {
				return errors.Errorf("failed to implement Serialized %v", dest)
			}
			err = rs.scripts["1_HGETALLXP"].Send(conn, rs.name, dest.PrimaryKey())
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
					dest.Deserialized(vv.([]byte))
					break
				}
			}
		default:
			panic("unsupported store type")
		}
	}

	return nil
}
