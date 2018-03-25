package goarmorapi

import (
	"fmt"
	"strings"

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

func (keyValues KVS) ResponseMessages() []*JSONMsg {
	if len(keyValues) == 0 {
		return nil
	}
	e := make([]*JSONMsg, 0, len(keyValues))
	for k, v := range keyValues {
		e = append(e, &JSONMsg{
			Code: uint64(MsgCodeDefaultDebug), Err: fmt.Errorf("%s:%s", k, v)})
	}
	return e
}

func newKVS(m map[string]string) KVS {
	if m == nil {
		m = make(map[string]string)
	}
	return KVS(m)
}
