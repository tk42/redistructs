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
		return errors.Wrapf(err, "failed to send WATCH %ss", rs.name)
	}
	_, err = conn.Do("WATCH", rs.key)
	if err != nil {
		return errors.Wrapf(err, "failed to send WATCH %ss", rs.key)
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
		if len(m.Serialized()) == 0 {
			return errors.Errorf("failed to implement Serialized %v", m)
		}
		err = rs.scripts["HSETXP"].Send(conn, rs.name, rs.getExpireArg(), m.PrimaryKey(), m.Serialized())
		if err != nil {
			return errors.Wrapf(err, "failed to send HSETXP %s", rs.key)
		}
		var args []string
		for k, v := range m.ScoreMap() {
			args = append(args, fmt.Sprint(v), fmt.Sprint(k))
		}
		err = rs.scripts["ZADDXP"].Send(conn, rs.key, rs.getExpireArg(), args)
		if err != nil {
			conn.Do("DISCARD")
			return errors.Wrapf(err, "failed to 2_ZADDXP %v", rs.key)
		}
	default:
		panic("unsupported store type")
	}

	vals, err := redigo.Values(conn.Do("EXEC"))
	if err != nil {
		return errors.Wrap(err, "faild to EXEC commands")
	}
	for _, r := range vals {
		if r.(int64) != 1 {
			return errors.Wrap(errors.New("return FAILED after EXEC commands"), fmt.Sprint(vals))
		}
	}
	// if vals[0] != "OK" {
	// 	return
	// }
	return nil
}
