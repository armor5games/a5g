package gameserverapi

const (
	ConfigKey ctxKey = iota
	LoggerWithoutUserKey
	DBReadPoolKey
	DBWritePoolKey
	DBReadStaticPoolKey
	StaticDataKey
	LoginUsersAvatarsCacheKey
	ShardUsersAvatarsCacheKey
)

type ctxKey int
