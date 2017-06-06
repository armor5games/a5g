package gameserverapi

const (
	CtxConfigKey CtxKey = iota
	CtxLoggerWithoutUserKey
	CtxDBReadPoolKey
	CtxDBWritePoolKey
	CtxDBReadStaticPoolKey
	CtxStaticCacheKey
	CtxUsersCacheKey
	CtxCurrentUserIDKey
)

type CtxKey int
