package gameserverapi

const (
	CtxConfigKey CtxKey = iota
	CtxLoggerWithoutUserKey
	CtxDBReadPoolKey
	CtxDBWritePoolKey
	CtxDBReadStaticPoolKey
	CtxStaticDataKey
	CtxAvatarLoginUsersKey
	CtxAvatarShardUsersCacheKey
)

type CtxKey int
