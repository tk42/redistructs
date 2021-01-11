package redistructs

import (
	"context"
	"fmt"
	"testing"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/tk42/redistructs/types"
)

func TestPut(t *testing.T) {
	now := time.Now()
	postStore := New(redisPool, *types.CreateConfig(), &Post{})

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
	if err != nil {
		panic(err)
	}

	conn := redisPool.Get()
	defer conn.Close()

	c, err := redigo.Int(conn.Do("HLEN", "*Post"))
	if err != nil {
		panic(err)
	}
	if c != 4 {
		panic(fmt.Errorf("expect: %v, got: %v", 4, c))
	}

	c, err = redigo.Int(conn.Do("ZCARD", "*Post/0"))
	if err != nil {
		panic(err)
	}
	if c != 2 {
		panic(fmt.Errorf("expect: %v, got: %v", 2, c))
	}
	c, _ = redigo.Int(conn.Do("ZCARD", "*Post/0.EXPIREAT"))
	if c != 2 {
		panic(fmt.Errorf("expect: %v, got: %v", 2, c))
	}
}

func TestGet(t *testing.T) {
	TestPut(t)

	postStore := New(redisPool, *types.CreateConfig(), &Post{})

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
}

func TestDelete(t *testing.T) {
	TestPut(t)

	postStore := New(redisPool, *types.CreateConfig(), &Post{})

	p := &Post{ID: 4}
	err := postStore.Delete(context.TODO(), p)
	if err != nil {
		panic(err)
	}

	conn := redisPool.Get()
	defer conn.Close()

	c, err := redigo.Int(conn.Do("HLEN", "*Post"))
	if err != nil {
		panic(err)
	}
	if c != 3 {
		panic(fmt.Errorf("expect: %v, got: %v", 3, c))
	}

	c, err = redigo.Int(conn.Do("ZCARD", "*Post/0"))
	if err != nil {
		panic(err)
	}
	if c != 2 {
		panic(fmt.Errorf("expect: %v, got: %v", 2, c))
	}
	c, _ = redigo.Int(conn.Do("ZCARD", "*Post/0.EXPIREAT"))
	if c != 2 {
		panic(fmt.Errorf("expect: %v, got: %v", 2, c))
	}
}
