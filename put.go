package redistructs

import (
	"context"
	"fmt"
	"reflect"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/tk42/redistructs/types"
)

func (rs *RedigoStructs) Put(ctx context.Context, src interface{}) error {
	conn, err := rs.pool.GetContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to acquire a connection")
	}
	defer conn.Close()

	rv := reflect.ValueOf(src)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			err = rs.set(conn, rv.Index(i))
			if err != nil {
				break
			}
		}
	} else {
		err = rs.set(conn, rv)
	}
	return err
}

func (rs *RedigoStructs) set(conn redigo.Conn, src reflect.Value) error {
	m, ok := src.Interface().(types.RediStruct)
	if !ok {
		return fmt.Errorf("failed to cast %v to types.RediStruct", src.Interface())
	}

	dbIdx, changed := rs.getDBIndex(rs.model)
	if changed {
		conn.Do("SELECT", dbIdx)
	}

	_, err := conn.Do("WATCH", rs.name)
	if err != nil {
		return errors.Wrapf(err, "failed to send WATCH %s", rs.name)
	}
	_, err = conn.Do("WATCH", rs.key)
	if err != nil {
		return errors.Wrapf(err, "failed to send WATCH %s", rs.key)
	}

	err = conn.Send("MULTI")
	if err != nil {
		return errors.Wrap(err, "faild to send MULTI command")
	}

	switch m.StoreType() {
	// case types.FlattenHash:
	// 	flat, err := flatten.Flatten(structs.Map(m), m.PrimaryKey(), flatten.SeparatorStyle{Middle: m.KeyDelimiter()})
	// 	if err != nil {
	// 		return errors.Errorf("failed to flatten %v", err)
	// 	}
	case types.Serialized:
		b, err := m.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal")
		}
		if len(b) == 0 {
			return errors.Errorf("failed to marshal due to empty %v", rs.model)
		}
		err = rs.scripts["HSETXP"].Send(conn, rs.name, rs.getExpireArg(), m.PrimaryKey(), b)
		if err != nil {
			return errors.Wrapf(err, "failed to send HSETXP %s", rs.key)
		}
		args := []interface{}{rs.key, rs.getExpireArg()}
		for k, v := range m.ScoreMap() {
			args = append(args, fmt.Sprint(v), fmt.Sprint(k))
		}
		err = rs.scripts["ZADDXP"].Send(conn, args...)
		if err != nil {
			conn.Do("DISCARD")
			return errors.Wrapf(err, "failed to ZADDXP %v", rs.key)
		}
	default:
		panic("unsupported store type")
	}

	vals, err := redigo.Ints(conn.Do("EXEC"))
	if err != nil {
		return errors.Wrap(err, "faild to EXEC commands")
	}
	for _, r := range vals {
		// the number of elements added to the sorted set
		if r == 0 {
			return errors.Wrap(errors.New("return FAILED after EXEC commands"), fmt.Sprint(vals))
		}
	}
	return nil
}
