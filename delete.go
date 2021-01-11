package redistructs

import (
	"context"
	"fmt"
	"reflect"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"github.com/tk42/redistructs/types"
)

func (rs *RedigoStructs) Delete(ctx context.Context, src interface{}) error {
	conn, err := rs.pool.GetContext(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to acquire a connection")
	}
	defer conn.Close()

	rv := reflect.ValueOf(src)
	if rv.Kind() == reflect.Slice {
		for i := 0; i < rv.Len(); i++ {
			err = rs.delete(conn, rv.Index(i))
			if err != nil {
				break
			}
		}
	} else {
		err = rs.delete(conn, rv)
	}
	return err
}

func (rs *RedigoStructs) delete(conn redigo.Conn, src reflect.Value) error {
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
		err = rs.scripts["1_HDELXP"].Send(conn, rs.name, m.PrimaryKey())
		if err != nil {
			return errors.Wrapf(err, "failed to send HDELXP %s", rs.key)
		}
		var args []string
		for k, _ := range m.ScoreMap() {
			args = append(args, fmt.Sprint(k))
		}
		err = rs.scripts["1_ZREM"].Send(conn, rs.key, args)
		if err != nil {
			conn.Do("DISCARD")
			return errors.Wrapf(err, "failed to 2_ZADDXP %v", rs.key)
		}
	default:
		panic("unsupported store type")
	}

	_, err = conn.Do("EXEC")
	if err != nil {
		return errors.Wrap(err, "failed to execute EXEC")
	}
	return nil
}
