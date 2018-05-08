package a5gfields

import (
	"fmt"
	"time"

	"github.com/armor5games/a5g/a5gvalues"
)

type Field interface {
	Key() string
	Value() string
}

func New(k string, v fmt.Stringer) *kvField { return &kvField{key: k, value: v} }

func (v *kvField) Key() string   { return v.key }
func (v *kvField) Value() string { return v.value.String() }

func Bytes(k string, v []byte) *kvField {
	return &kvField{key: k, value: a5gvalues.Bytes(v)}
}

func Duration(k string, v time.Duration) *kvField {
	return &kvField{key: k, value: a5gvalues.Duration(v)}
}

func EmptyInterface(k string, v interface{}) *kvField {
	return &kvField{key: k, value: a5gvalues.EmptyInterface(v)}
}

func Float64(k string, v float64) *kvField {
	return &kvField{key: k, value: a5gvalues.Float64(v)}
}

func Int(k string, v int) *kvField {
	return &kvField{key: k, value: a5gvalues.Int(v)}
}

func Int64(k string, v int64) *kvField {
	return &kvField{key: k, value: a5gvalues.Int64(v)}
}

func String(k string, v string) *kvField {
	return &kvField{key: k, value: a5gvalues.String(v)}
}

// kvField is an key-value pair (kv)
type kvField struct {
	key   string
	value fmt.Stringer
}
