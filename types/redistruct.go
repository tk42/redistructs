package types

import "time"

// RediStruct is an interface for hundling redistructs
type RediStruct interface {
	// Permission() Permission
	StoreType() StoreType
	PrimaryKey() string
	KeyDelimiter() string
	ScoreMap() map[string]interface{}
	Expire() time.Duration
	DatabaseIdx() int
	Serialized() []byte
	Deserialized([]byte)
}
