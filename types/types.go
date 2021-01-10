package types

type ClientType string

// type Permission string

type StoreType string

const (
	// ReadOnly     = Permission("read-only")
	// ReadAndWrite = Permission("read&write")

	FlattenHash = StoreType("flatten-hash")
	Serialized  = StoreType("serialized")
	// NestedHash  = StoreType("nest-hash")

	Redigo = ClientType("redigo")
)
