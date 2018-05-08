package a5gkv

import (
	"fmt"
	"strings"

	"github.com/armor5games/a5g/a5gfields"
	"github.com/armor5games/a5g/a5gvalues"
	"github.com/pkg/errors"
)

// KV is an key-values
type KV map[string]fmt.Stringer

func (m KV) Error() string { return m.String() }

func New() KV { return newKV() }

func NewByMapString(m map[string]string) KV {
	m2 := newKV()
	for k, v := range m {
		m2[k] = a5gvalues.String(v)
	}
	return m2
}

func (m KV) String() string {
	if len(m) == 0 {
		return ""
	}
	a := make([]string, 0, len(m))
	for k, v := range m {
		a = append(a, fmt.Sprintf("%s:%v", k, v.String()))
	}
	return fmt.Sprintf("kv:[%s]", strings.Join(a, " "))
}

func (m KV) Err() error {
	if len(m) == 0 {
		return nil
	}
	return errors.New(m.String())
}

func (m KV) Fields() []a5gfields.Field {
	if len(m) == 0 {
		return nil
	}
	a := make([]a5gfields.Field, 0, len(m))
	for k, v := range m {
		a = append(a, a5gfields.New(k, v))
	}
	return a
}

func (m KV) Merge(m2 KV) { m.merge(m2) }

func (m KV) Copy() KV {
	var m2 = New()
	for k, v := range m {
		m2[k] = v
	}
	return m2
}

// func (m KV) ResponseMessages() []*a5gapi.APIErr {
// 	if len(m) == 0 {
// 		return nil
// 	}
// 	e := make([]*a5gapi.APIErr, 0, len(m))
// 	for k, v := range m {
// 		e = append(e, &a5gapi.APIErr{
// 			Code: uint64(a5gapi.ErrCodeDefaultDebug),
// 			Err:  fmt.Errorf("%s:%s", k, v.String())})
// 	}
// 	return e
// }

func newKV(a ...map[string]fmt.Stringer) KV {
	if len(a) == 0 {
		return KV(make(map[string]fmt.Stringer))
	}
	var (
		m  = KV(a[0])
		a2 []KV
	)
	for i := 1; i < len(a); i++ {
		a2 = append(a2, KV(a[i]))
	}
	m.merge(a2...)
	return m
}

func (m KV) merge(a ...KV) {
	for _, m2 := range a {
		for k, v := range m2 {
			m[k] = v
		}
	}
}
