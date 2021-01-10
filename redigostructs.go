package redistructs

import (
	"context"
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
	RediStructs
	config   types.Config
	model    types.RediStruct
	dbIdx    int
	name     string
	key      string
	expire   time.Duration
	expireAt time.Time
	pool     Pool
	scripts  map[string]*redigo.Script
}

func NewRedigoStructs(config types.Config, model types.RediStruct) RediStructs {
	redigoStructs := RedigoStructs{
		config: config,
		model:  model,
		dbIdx:  0,
		name:   types.GetName(model),
		key:    types.GetName(model) + model.KeyDelimiter() + model.PrimaryKey(),
		pool: &redigo.Pool{
			MaxIdle:     config.MaxIdle,
			MaxActive:   config.MaxActive,
			Wait:        false, // true: blocking until the number of connections is under MaxActive
			IdleTimeout: config.IdleTimeout,
			Dial: func() (redigo.Conn, error) {
				return redigo.Dial("tcp", config.Address)
			},
		},
		scripts: loadScripts(),
	}
	redigoStructs.setExpire(model)
	return redigoStructs
}

func (rs RedigoStructs) setExpire(model types.RediStruct) {
	e := model.Expire()
	switch e.(type) {
	case time.Time:
		rs.expireAt = e.(time.Time)
	case time.Duration:
		rs.expire = e.(time.Duration)
	default:
		panic("Invalid model.Expire()")
	}
}

func (rs RedigoStructs) getExpire() time.Duration {
	if rs.expireAt.IsZero() {
		return rs.expire
	} else {
		return rs.expireAt.Sub(time.Now())
	}
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

func (rs RedigoStructs) getDBIndex(model types.RediStruct) (int, bool) {
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

// func (rs RedigoStructs) Count(ctx context.Context, mods ...rq.Modifier) (int, error) {
// 	return 0, nil
// }

// func (rs RedigoStructs) Map() map[string]types.RediStruct {
// }

// func (rs RedigoStructs) Values() []types.RediStruct {
// }

// func (rs RedigoStructs) Names() []string {
// }

// func (rs RedigoStructs) IsZero() bool {
// }

// func (rs RedigoStructs) HasZero() bool {
// }
