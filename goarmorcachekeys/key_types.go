package goarmorcachekeys

import (
	"bytes"
	"fmt"
)

type KeyType int

func (t KeyType) Sprint(a ...interface{}) string {
	if len(a) == 0 {
		return fmt.Sprintf("%d", t)
	}

	var b bytes.Buffer

	for _, v := range a {
		b.WriteString(fmt.Sprint(v))
	}

	return fmt.Sprintf("%d%s", t, b.String())
}
