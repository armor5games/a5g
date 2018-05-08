package a5gapi

import (
	"fmt"
	"strings"

	"github.com/armor5games/a5g/a5gfields"
	"github.com/armor5games/a5g/a5gvalues"
	"github.com/pkg/errors"
)

// KVS is an key and values of type string
type KVS map[string]string

type KVSValueBoolean string

func (keyValues KVS) Error() string { return keyValues.String() }

func NewKVS() KVS { return newKVS(nil) }

const (
	KVSValBoolTrue  KVSValueBoolean = "true"
	KVSValBoolFalse KVSValueBoolean = "false"
)

func (v KVSValueBoolean) String() string { return string(v) }

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

func (keyValues KVS) Err() error {
	if len(keyValues) == 0 {
		return nil
	}
	return errors.New(keyValues.String())
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

func (m KVS) Fields() []a5gfields.Field {
	if len(m) == 0 {
		return nil
	}
	a := make([]a5gfields.Field, 0, len(m))
	for k, v := range m {
		a = append(a, a5gfields.New(k, a5gvalues.String(v)))
	}
	return a
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

func (keyValues KVS) ResponseMessages() []*APIErr {
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

func newKVS(m map[string]string) KVS {
	if m == nil {
		m = make(map[string]string)
	}
	return KVS(m)
}
