package gameserverapi

const (
	CtxConfigKey CtxKey = iota
	CtxLoggerWithoutUserKey
	CtxDBReadPoolKey
	CtxDBWritePoolKey
	CtxDBReadStaticPoolKey
	CtxStaticCacheKey
	CtxUsersLCacheKey
	CtxUsersSCacheKey
	CtxCurrentUserIDKey
)

type CtxKey int
