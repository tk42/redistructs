package types

import "time"

type Config struct {
	client            ClientType
	DatabaseIdx       int
	PrimaryKey        string
	SuffixScoreSetKey string
	expire            time.Duration
	expireAt          time.Time
}

type ConfigOption func(*Config)

func DatabaseIdx(databaseIdx int) ConfigOption {
	return func(op *Config) {
		op.DatabaseIdx = databaseIdx
	}
}

func PrimaryKey(PrimaryKey string) ConfigOption {
	return func(op *Config) {
		op.PrimaryKey = PrimaryKey
	}
}

func SuffixScoreSetKey(suffixScoreSetKey string) ConfigOption {
	return func(op *Config) {
		op.SuffixScoreSetKey = suffixScoreSetKey
	}
}

func CreateConfig(ops ...ConfigOption) *Config {
	config := Config{
		client:            Redigo,
		DatabaseIdx:       0,
		PrimaryKey:        "",
		SuffixScoreSetKey: "/SCORE",
	}
	for _, option := range ops {
		option(&config)
	}
	return &config
}
