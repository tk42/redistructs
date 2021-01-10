module github.com/tk42/redistructs

go 1.15

replace (
	github.com/tk42/redistructs/goredis => ./goredis
	github.com/tk42/redistructs/redigo => ./redigo
	github.com/tk42/redistructs/structs => ./structs
	github.com/tk42/redistructs/types => ./types
)

require (
	github.com/fatih/structs v1.1.0
	github.com/garyburd/redigo v1.6.2
	github.com/go-errors/errors v1.1.1 // indirect
	github.com/go-redis/redis/v8 v8.4.4
	github.com/go-sql-driver/mysql v1.5.0 // indirect
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/izumin5210/ro v0.4.0
	github.com/lib/pq v1.9.0 // indirect
	github.com/ory-am/common v0.4.0 // indirect
	github.com/ory/dockertest v3.3.2+incompatible
	github.com/ory/dockertest/v3 v3.6.3
	github.com/pborman/uuid v1.2.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/tk42/redis-lua-scripts v0.0.0-20210106105717-72cb85d4286d
)
