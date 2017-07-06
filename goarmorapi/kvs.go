package goarmorapi

import "fmt"

// KVS is an key and values of type string
type KVS map[string]string

func NewKVS() KVS {
	return newKVS(nil)
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

func (keyValues KVS) ResponseErrors() []*ErrorJSON {
	if len(keyValues) == 0 {
		return nil
	}

	e := make([]*ErrorJSON, 0, len(keyValues))

	for k, v := range keyValues {
		e = append(e, &ErrorJSON{
			Code: KVAPIErrorCode, Error: fmt.Errorf("%s:%s", k, v)})
	}

	return e
}

func newKVS(m map[string]string) KVS {
	if m == nil {
		m = make(map[string]string)
	}

	return KVS(m)
}
