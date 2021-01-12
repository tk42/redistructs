package types

import (
	"reflect"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/izumin5210/ro/rq"
	"github.com/pkg/errors"
)

func GetName(myvar interface{}) string {
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func strctVal(s interface{}) reflect.Value {
	v := reflect.ValueOf(s)

	// if pointer get the underlying element≤
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("not struct")
	}

	return v
}

func selectKeys(conn redigo.Conn, prefix string, mods []rq.Modifier) ([]string, error) {
	q := rq.List(mods...)
	q.Key.Prefix = prefix
	cmd, err := q.Build()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	keys, err := redigo.Strings(conn.Do(cmd.Name, cmd.Args...))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return keys, nil
}
