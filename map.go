package redistructs

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/tk42/redistructs/types"
)

func (rs *RedigoStructs) Map(ctx context.Context) (map[string]types.RediStruct, error) {
	conn, err := rs.pool.GetContext(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to acquire a connection")
	}
	defer conn.Close()

	dbIdx, changed := rs.getDBIndex(rs.model)
	if changed {
		conn.Do("SELECT", dbIdx)
	}

	switch rs.model.StoreType() {
	case types.Serialized:
		if len(rs.model.Serialized()) == 0 {
			return nil, errors.Errorf("failed to implement Serialized %v", rs.model)
		}
		err = rs.scripts["HGETALLXP"].Send(conn, rs.name, rs.model.PrimaryKey())
		if err != nil {
			return nil, errors.Wrapf(err, "failed to send HGETALLXP %s", rs.key)
		}
	default:
		panic("unsupported store type")
	}

	err = conn.Flush()
	if err != nil {
		return nil, errors.Wrap(err, "faild to flush commands")
	}

	v, err := redigo.Values(conn.Receive())
	if err != nil {
		return nil, errors.Wrap(err, "faild to receive or cast redis command result")
	}

	result := make(map[string]types.RediStruct)

	switch rs.model.StoreType() {
	case types.Serialized:
		var key string
		for j, vv := range v {
			if j%2 == 0 {
				key = string(vv.([]byte))
			} else {
				s := rs.clone().(types.RediStruct)
				s.Deserialized(vv.([]byte))
				result[key] = s
			}
		}
	default:
		panic("unsupported store type")
	}

	return result, nil
}
