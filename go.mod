module github.com/tk42/redistructs

go 1.15

replace (
	github.com/tk42/redistructs/goredis => ./goredis
	github.com/tk42/redistructs/redigo => ./redigo
	github.com/tk42/redistructs/structs => ./structs
	github.com/tk42/redistructs/types => ./types
)

require (
	github.com/gomodule/redigo v2.0.0+incompatible
	github.com/google/go-cmp v0.5.4 // indirect
	github.com/izumin5210/ro v0.4.0
	github.com/lib/pq v1.9.0 // indirect
	github.com/ory/dockertest/v3 v3.6.3
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.6.1 // indirect
	github.com/tk42/redis-lua-scripts v0.0.0-20210110151629-279682391ee1
	golang.org/x/net v0.0.0-20201202161906-c7110b5ffcbb // indirect
)
