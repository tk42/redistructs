package redistructs

import (
	"context"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/tk42/redistructs/types"
)

func (rs *RedigoStructs) Names(ctx context.Context) ([]string, error) {
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
		err = rs.scripts["HKEYSXP"].Send(conn, rs.name)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to send HKEYS %s", rs.key)
		}
	default:
		panic("unsupported store type")
	}

	err = conn.Flush()
	if err != nil {
		return nil, errors.Wrap(err, "faild to flush commands")
	}

	v, err := redigo.Strings(conn.Receive())
	if err != nil {
		return nil, errors.Wrap(err, "faild to receive or cast redis command result")
	}

	var result []string

	switch rs.model.StoreType() {
	case types.Serialized:
		for _, vv := range v {
			result = append(result, vv)
		}
	default:
		panic("unsupported store type")
	}

	return result, nil
}
