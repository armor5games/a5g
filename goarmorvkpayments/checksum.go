package goarmorvkpayments

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"sort"
)

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
