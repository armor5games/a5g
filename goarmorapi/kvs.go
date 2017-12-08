package goarmorapi

import (
	"fmt"
	"strings"
)

// KVS is an key and values of type string
type KVS map[string]string

func NewKVS() KVS {
	return newKVS(nil)
}

type KVSValueBoolean string

const (
	KVSValBoolTrue  KVSValueBoolean = "true"
	KVSValBoolFalse KVSValueBoolean = "false"
)

func (v KVSValueBoolean) String() string {
	return string(v)
}

func (keyValues KVS) String() string {
	if len(keyValues) == 0 {
		return ""
	}

	kv := make([]string, 0, len(keyValues))

	for k, v := range keyValues {
		kv = append(kv, fmt.Sprintf("%s:%s", k, v))
	}

	return fmt.Sprintf("kvs:[%s]", strings.Join(kv, " "))
}

func (keyValues KVS) KV() KV {
	if len(keyValues) == 0 {
		return nil
	}

	kv := make(KV)

	for k, v := range keyValues {
		kv[k] = v
	}

	return kv
}

func (keyValues KVS) Merge(newKeyValues KVS) {
	for k, v := range newKeyValues {
		keyValues[k] = v
	}
}

func (keyValues KVS) Copy() KVS {
	var newKeyValues = NewKVS()

	for k, v := range keyValues {
		newKeyValues[k] = v
	}

	return newKeyValues
}

func (keyValues KVS) ResponseErrors() []*ErrorJSON {
	if len(keyValues) == 0 {
		return nil
	}

	e := make([]*ErrorJSON, 0, len(keyValues))

	for k, v := range keyValues {
		e = append(e, &ErrorJSON{
			Code: uint64(ErrCodeDefautlDebug), Err: fmt.Errorf("%s:%s", k, v)})
	}

	return e
}

func newKVS(m map[string]string) KVS {
	if m == nil {
		m = make(map[string]string)
	}

	return KVS(m)
}
