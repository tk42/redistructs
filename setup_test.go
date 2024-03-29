package redistructs

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/tk42/redistructs/types"
)

var (
	redisPool *types.Pool
	TTL       = int64(2)
)

type Post struct {
	ID        int64  `redis:"id"`
	UserID    uint64 `redis:"user_id"`
	Title     string `redis:"title"`
	Body      string `redis:"body"`
	CreatedAt int64  `redis:"created_at"`
}

func (p *Post) StoreType() types.StoreType {
	return types.Serialized
}

func (p *Post) PrimaryKey() string {
	return fmt.Sprint(p.ID)
}

func (p *Post) KeyDelimiter() string {
	return "/"
}

func (p *Post) ScoreMap() map[string]interface{} {
	return map[string]interface{}{
		"id":     p.ID,
		"recent": p.CreatedAt,
	}
}

func (p *Post) Expire() time.Duration {
	return time.Duration(time.Second)
}

func (p *Post) DatabaseIdx() int {
	return 0
}

func (p *Post) Marshal() ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	err := gob.NewEncoder(buf).Encode(&p)
	if err != nil {
		panic("Failed to Serialized")
	}
	return buf.Bytes(), err
}

func (p *Post) Unmarshal(b []byte) error {
	return gob.NewDecoder(bytes.NewBuffer(b)).Decode(&p)
}

func TestMain(m *testing.M) {
	redisPool = types.MustCreate()

	code := m.Run()

	redisPool.MustClose()

	os.Exit(code)
}
