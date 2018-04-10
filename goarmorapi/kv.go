package goarmorapi

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// KV is an key-values
type KV map[string]interface{}

func (keyValues KV) Error() string { return keyValues.String() }

func NewKV() KV { return newKV(nil) }

func (keyValues KV) String() string {
	if len(keyValues) == 0 {
		return ""
	}
	kv := make([]string, 0, len(keyValues))
	for k, v := range keyValues {
		kv = append(kv, fmt.Sprintf("%s:%s", k, v))
	}
	return fmt.Sprintf("kv:[%s]", strings.Join(kv, " "))
}

func (keyValues KV) Err() error {
	if len(keyValues) == 0 {
		return nil
	}
	return errors.New(keyValues.String())
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

func (keyValues KV) Copy() KV {
	var newKeyValues = NewKV()
	for k, v := range keyValues {
		newKeyValues[k] = v
	}
	return newKeyValues
}

func (keyValues KV) ResponseMessages() []*APIErr {
	if len(keyValues) == 0 {
		return nil
	}
	e := make([]*APIErr, 0, len(keyValues))
	for k, v := range keyValues {
		e = append(e, &APIErr{
			Code: uint64(ErrCodeDefaultDebug), Err: fmt.Errorf("%s:%s", k, v)})
	}
	return e
}

func newKV(m map[string]interface{}) KV {
	if m == nil {
		m = make(map[string]interface{})
	}
	return KV(m)
}
