package goarmorvkpayments

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"sort"
)

type KV map[string]string

func NewChecksum(vkKV KV, secretKey string) (string, error) {
	var kvKeys []string

	for k := range vkKV {
		kvKeys = append(kvKeys, k)
	}

	sort.Strings(kvKeys)

	var b bytes.Buffer

	for _, k := range kvKeys {
		b.WriteString(fmt.Sprintf("%s=%s", k, vkKV[k]))
	}

	h := md5.New()

	_, err := io.WriteString(h, fmt.Sprintf("%s%s", b.String(), secretKey))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
