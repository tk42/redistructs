package types

// RediStruct is an interface for hundling redistructs
type RediStruct interface {
	// Permission() Permission
	StoreType() StoreType
	PrimaryKey() string
	KeyDelimiter() string
	ScoreMap() map[string]interface{}
	Expire() interface{} // time.Duration || time.Time
	DatabaseIdx() int
	Serialized() []byte
	Deserialized([]byte)
}
