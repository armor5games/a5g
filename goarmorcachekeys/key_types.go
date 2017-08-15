package goarmorcachekeys

import (
	"bytes"
	"fmt"
)

type KeyType int

func (t KeyType) Sprint(a ...interface{}) string {
	var b bytes.Buffer

	for _, v := range a {
		b.WriteString(fmt.Sprint(v))
	}

	return fmt.Sprintf("%d%v", t, b.String())
}
