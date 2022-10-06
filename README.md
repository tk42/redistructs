# redistructs
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Boilerplate codes to store struct as hashset or serialized object into Redis

## Features
 - RediStructs has an interface compatible with fatih/structs which is optimizely queried.

 - RediStore is an object for operations to Model and Redistructs inspired by izumin5210/ro

 - Set EXPIRE/EXPIREAT

 - Set a prefix string

 - Atomic read and write

 - Store recursive redistructs if the struct is nested or serialized object (by gob)

## Examples
### Store a struct
```go
redisPool := &redis.Pool{
    Dial: func() (redis.Conn, error) {
        return redis.DialURL(fmt.Sprintf("redis://localhost:%s", p.dockerRes.GetPort("6379/tcp")))
    }
}
postStore := New(redisPool, types.CreateConfig(), &Post{})

err := postStore.Put(context.TODO(), []*Post{
    {
        ID:        1,
        UserID:    1,
        Title:     "post 1",
        Body:      "This is a post 1",
        CreatedAt: now.UnixNano(),
    },
    {
        ID:        2,
        UserID:    2,
        Title:     "post 2",
        Body:      "This is a post 2",
        CreatedAt: now.Add(-24 * 60 * 60 * time.Second).UnixNano(),
    },
    {
        ID:        3,
        UserID:    1,
        Title:     "post 3",
        Body:      "This is a post 3",
        CreatedAt: now.Add(24 * 60 * 60 * time.Second).UnixNano(),
    },
    {
        ID:        4,
        UserID:    1,
        Title:     "post 4",
        Body:      "This is a post 4",
        CreatedAt: now.Add(-24 * 60 * 60 * time.Second).UnixNano(),
    },
})
```

### Get a struct
```go
postStore := New(redisPool, types.CreateConfig(), &Post{})

p := &Post{ID: 4}
err := postStore.Get(context.TODO(), p)
if err != nil {
    panic(err)
}

if p.UserID != 1 {
    panic(fmt.Errorf("expect: %v, got: %v", 1, p.UserID))
}

if p.Title != "post 4" {
    panic(fmt.Errorf("expect: %v, got: %v", "post 4", p.Title))
}
```

### Delete a struct
```go
postStore := New(redisPool, types.CreateConfig(), &Post{})

p := &Post{ID: 4}
err := postStore.Delete(context.TODO(), p)
if err != nil {
    panic(err)
}
```