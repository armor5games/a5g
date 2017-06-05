package gameserverapi

const (
	CtxConfigKey CtxKey = iota
	CtxLoggerWithoutUserKey
	CtxDBReadPoolKey
	CtxDBWritePoolKey
	CtxDBReadStaticPoolKey
	CtxStaticDataKey
	CtxLoginUsersAvatarsCacheKey
	CtxShardUsersAvatarsCacheKey
)

type CtxKey int
