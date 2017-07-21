package goarmorcachekeys

import (
	"bytes"
	"fmt"
)

type KeyType int

func (t KeyType) Sprint(a ...interface{}) string {
	var b bytes.Buffer

	for _, i := range a {
		b.WriteString(fmt.Sprint(i))
	}

	return fmt.Sprintf("%d%v", t, b.String())
}
