package a5gvalues

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func Bytes(v []byte) *kvValue {
	return &kvValue{typ: bytesType, valueEmptyInterface: v}
}

func Duration(v time.Duration) *kvValue {
	return &kvValue{typ: durationType, valueInt64: int64(v)}
}

func EmptyInterface(v interface{}) *kvValue {
	return &kvValue{typ: emptyInterfaceType, valueEmptyInterface: v}
}

func Float64(v float64) *kvValue {
	return &kvValue{typ: float64Type, valueInt64: int64(math.Float64bits(v))}
}

func Int(v int) *kvValue {
	return &kvValue{typ: intType, valueInt64: int64(v)}
}

func String(v string) *kvValue {
	return &kvValue{typ: stringType, valueString: v}
}

func (v *kvValue) String() string {
	switch v.typ {
	case bytesType:
		return fmt.Sprintf("%+v", v.valueEmptyInterface.([]byte))
	case stringType:
		return v.valueString
	case durationType:
		return fmt.Sprintf("%+v", time.Duration(v.valueInt64))
	case float64Type:
		return fmt.Sprintf("%+v", math.Float64frombits(uint64(v.valueInt64)))
	case intType:
		return strconv.FormatInt(v.valueInt64, 10)
	case emptyInterfaceType:
		return fmt.Sprintf("%+v", v.valueEmptyInterface)
	}
	panic(fmt.Sprintf("unknown value type: %v", v.typ))
}

// kvValue is an key-value's value
type kvValue struct {
	typ                 valueType
	valueInt64          int64
	valueString         string
	valueEmptyInterface interface{}
}

type valueType uint8

const (
	// unknownType is the default value type. Attempting to add it to an encoder will panic.
	unknownType valueType = iota
	bytesType
	durationType
	float64Type
	intType
	stringType
	// EmptyInterfaceType indicates that the field carries an interface{}, which should
	// be serialized using reflection.
	emptyInterfaceType
)
