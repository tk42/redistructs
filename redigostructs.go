package redistructs

import (
	"context"
	"fmt"
	"reflect"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/tk42/redistructs/types"

	rls "github.com/tk42/redis-lua-scripts"
)

// Pool is a pool of redis connections.
type Pool interface {
	GetContext(context.Context) (redigo.Conn, error)
}

type RedigoStructs struct {
	config  types.Config
	model   types.RediStruct
	dbIdx   int
	name    string
	key     string
	expire  time.Duration
	pool    Pool
	scripts map[string]*redigo.Script
}

func NewRedigoStructs(pool Pool, config types.Config, model types.RediStruct) RediStructs {
	redigoStructs := RedigoStructs{
		config:  config,
		model:   model,
		dbIdx:   0,
		name:    types.GetName(model),
		key:     types.GetName(model) + model.KeyDelimiter() + model.PrimaryKey(),
		pool:    pool,
		scripts: loadScripts(),
	}
	redigoStructs.setExpire(model)
	return &redigoStructs
}

func (rs *RedigoStructs) setExpire(model types.RediStruct) {
	rs.expire = model.Expire()
}

func (rs *RedigoStructs) getExpireArg() string {
	return fmt.Sprint(float64(rs.expire) / float64(time.Second))
}

func loadScripts() map[string]*redigo.Script {
	scripts := make(map[string]*redigo.Script)
	for _, group := range []string{"HASHES_XP", "SORTED_SETS_XP"} {
		script, err := rls.GetAllScripts(group)
		if err != nil {
			panic(err)
		}
		for k, v := range script {
			scripts[k] = v
		}
	}
	return scripts
}

func (rs *RedigoStructs) getDBIndex(model types.RediStruct) (int, bool) {
	dbIdx := 0
	if rs.config.DatabaseIdx < 0 && model.DatabaseIdx() < 0 {
		panic("Invalid database index")
	} else if rs.config.DatabaseIdx < 0 {
		dbIdx = model.DatabaseIdx()
	} else if model.DatabaseIdx() < 0 {
		dbIdx = rs.config.DatabaseIdx
	}
	if dbIdx != rs.dbIdx {
		rs.dbIdx = dbIdx
		return dbIdx, true
	} else {
		return dbIdx, false
	}
}

func (e *RedigoStructs) clone() interface{} {
	nInter := reflect.New(reflect.TypeOf(e.model).Elem())

	val := reflect.ValueOf(e.model).Elem()
	nVal := nInter.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		nvField.Set(val.Field(i))
	}

	return nInter.Interface()
}
