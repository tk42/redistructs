package types

import (
	"time"
)

type Config struct {
	client            ClientType
	Address           string
	MaxIdle           int
	MaxActive         int
	IdleTimeout       time.Duration
	DatabaseIdx       int
	PrimaryKey        string
	SuffixScoreSetKey string
}

type ConfigOption func(*Config)

func Address(address string) ConfigOption {
	return func(op *Config) {
		op.Address = address
	}
}

func MaxIdle(maxIdle int) ConfigOption {
	return func(op *Config) {
		if maxIdle < 0 {
			return
		}
		op.MaxIdle = maxIdle
	}
}

func MaxActive(maxActive int) ConfigOption {
	return func(op *Config) {
		if maxActive < 0 {
			return
		}
		op.MaxActive = maxActive
	}
}

func IdleTimeout(idleTimeout time.Duration) ConfigOption {
	return func(op *Config) {
		op.IdleTimeout = idleTimeout
	}
}

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
		Address:           "localhost:6379",
		MaxIdle:           3,
		MaxActive:         0,
		IdleTimeout:       240 * time.Second,
		DatabaseIdx:       0,
		PrimaryKey:        "",
		SuffixScoreSetKey: "/SCORE",
	}
	for _, option := range ops {
		option(&config)
	}
	return &config
}
