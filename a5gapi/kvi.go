package a5gapi

import (
	"fmt"
	"strings"

	"github.com/armor5games/a5g/a5gfields"
	"github.com/armor5games/a5g/a5gvalues"
	"github.com/pkg/errors"
)

// KV is an key-values
type KV map[string]interface{}

func (m KV) Error() string { return m.String() }

func NewKV() KV { return newKV(nil) }

func (m KV) String() string {
	if len(m) == 0 {
		return ""
	}
	kv := make([]string, 0, len(m))
	for k, v := range m {
		kv = append(kv, fmt.Sprintf("%s:%s", k, v))
	}
	return fmt.Sprintf("kv:[%s]", strings.Join(kv, " "))
}

func (m KV) Err() error {
	if len(m) == 0 {
		return nil
	}
	return errors.New(m.String())
}

func (m KV) KVS() KVS {
	if len(m) == 0 {
		return nil
	}
	kvs := make(KVS)
	for k, v := range m {
		kvs[k] = fmt.Sprint(v)
	}
	return kvs
}

func (m KV) Fields() []a5gfields.Field {
	if len(m) == 0 {
		return nil
	}
	a := make([]a5gfields.Field, 0, len(m))
	for k, v := range m {
		a = append(a, a5gfields.New(k, a5gvalues.EmptyInterface(v)))
	}
	return a
}

func (m KV) Merge(m2 KV) {
	for k, v := range m2 {
		m[k] = v
	}
}

func (m KV) Copy() KV {
	var m2 = NewKV()
	for k, v := range m {
		m2[k] = v
	}
	return m2
}

func (m KV) ResponseMessages() []*APIErr {
	if len(m) == 0 {
		return nil
	}
	a := make([]*APIErr, 0, len(m))
	for k, v := range m {
		a = append(a, &APIErr{
			Code: uint64(ErrCodeDefaultDebug), Err: fmt.Errorf("%s:%s", k, v)})
	}
	return a
}

func newKV(m map[string]interface{}) KV {
	if m == nil {
		m = make(map[string]interface{})
	}
	return KV(m)
}
