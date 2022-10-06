package redistructs

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/tk42/redistructs/types"
)

func (rs *RedigoStructs) Values(ctx context.Context) ([]types.RediStruct, error) {
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
		b, err := rs.model.Marshal()
		if err != nil {
			return nil, errors.Wrap(err, "failed to marshal")
		}
		if len(b) == 0 {
			return nil, errors.Errorf("failed to marshal due to empty %v", rs.model)
		}
		err = rs.scripts["HVALSXP"].Send(conn, rs.name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to send HVALS %s", rs.key)
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

	var result []types.RediStruct

	switch rs.model.StoreType() {
	case types.Serialized:
		for _, vv := range v {
			s := rs.clone().(types.RediStruct)
			if err := s.Unmarshal(vv.([]byte)); err != nil {
				return nil, errors.Wrap(err, "failed to unmarshal value")
			}
			result = append(result, s)
		}
	default:
		panic("unsupported store type")
	}

	return result, nil
}
