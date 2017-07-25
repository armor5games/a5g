package goarmorvkpayments

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"net/url"
	"sort"

	"github.com/armor5games/goarmor/goarmorapi"
)

// VKAPIKV vk payment's key values
type VKAPIKV map[string]string

func NewVKAPIKVByURLValues(keyValues url.Values) VKAPIKV {
	kv := make(VKAPIKV)

	if keyValues == nil {
		return kv
	}

	for k, v := range keyValues {
		if len(v) > 0 {
			kv[k] = v[0]
		}
	}

	return kv
}

func (keyValues VKAPIKV) KV() goarmorapi.KV {
	kv := make(goarmorapi.KV)

	if len(keyValues) == 0 {
		return kv
	}

	for k, v := range keyValues {
		kv[k] = v
	}

	return kv
}

func (keyValues VKAPIKV) KVS() goarmorapi.KVS {
	return goarmorapi.KVS(keyValues)
}

func (keyValues VKAPIKV) Checksum(secretKey string) (string, error) {
	var kvKeys []string

	for k := range keyValues {
		kvKeys = append(kvKeys, k)
	}

	sort.Strings(kvKeys)

	var b bytes.Buffer

	for _, k := range kvKeys {
		if k == "sig" {
			continue
		}

		b.WriteString(fmt.Sprintf("%s=%s", k, keyValues[k]))
	}

	h := md5.New()

	_, err := io.WriteString(h, fmt.Sprintf("%s%s", b.String(), secretKey))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
