package goarmorapi

import (
	"fmt"
	"strings"
)

// KV is an key-values
type KV map[string]interface{}

func NewKV() KV {
	return newKV(nil)
}

func (keyValues KV) KVS() KVS {
	if len(keyValues) == 0 {
		return nil
	}

	kvs := make(KVS)

	for k, v := range keyValues {
		kvs[k] = fmt.Sprint(v)
	}

	return kvs
}

func (keyValues KV) Merge(newKeyValues KV) {
	for k, v := range newKeyValues {
		keyValues[k] = v
	}
}

func (keyValues KV) String() string {
	if len(keyValues) == 0 {
		return ""
	}

	kv := make([]string, 0, len(keyValues))

	for k, v := range keyValues {
		kv = append(kv, fmt.Sprintf("%s:%s", k, v))
	}

	return fmt.Sprintf("kvs:[%s]", strings.Join(kv, " "))
}

func (keyValues KV) ResponseErrors() []*ErrorJSON {
	if len(keyValues) == 0 {
		return nil
	}

	e := make([]*ErrorJSON, 0, len(keyValues))

	for k, v := range keyValues {
		e = append(e, &ErrorJSON{
			Code: uint64(ErrCodeDefautlDebug), Error: fmt.Errorf("%s:%s", k, v)})
	}

	return e
}

func newKV(m map[string]interface{}) KV {
	if m == nil {
		m = make(map[string]interface{})
	}

	return KV(m)
}
