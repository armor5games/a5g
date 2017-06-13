package goarmorapi

const (
	CtxKeyConfig CtxKey = iota
	CtxKeyLogger
	CtxKeyRandSrc
)

type CtxKey int
