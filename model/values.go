package model

import (
	"strconv"

	"github.com/pkg/errors"
)

func NewVariable(s string) Variable {
	return &variable{
		nameComponent: nameComponent(s),
	}
}

func (v variable) Kind() Kind {
	return VariableKind
}

func (v variable) Value() interface{} {
	return v.Name()
}

func ParseIntValue(s string) (Value, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse int`)
	}
	return NewIntValue(int(v)), nil
}

func NewIntValue(v int) Value {
	return &intValue{
		value: v,
	}
}

func (v intValue) Value() interface{} {
	return v.value
}

func (v intValue) Kind() Kind {
	return IntKind
}

func NewFloatValue(s string) (Value, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse float`)
	}
	return &floatValue{
		value: v,
	}, nil
}

func (v floatValue) Kind() Kind {
	return FloatKind
}

func (v floatValue) Value() interface{} {
	return v.value
}

func NewStringValue(s string) Value {
	return &stringValue{
		value: s,
	}
}

func (v stringValue) Kind() Kind {
	return StringKind
}

func (v stringValue) Value() interface{} {
	return v.value
}

func NewBoolValue(s string) (Value, error) {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse bool`)
	}
	return &boolValue{
		value: v,
	}, nil
}

func (v boolValue) Value() interface{} {
	return v.value
}

func (v boolValue) Kind() Kind {
	return BooleanKind
}

var nullv nullValue
func NullValue() Value {
	return &nullv
}

func (v nullValue) Value() interface{} {
	return nil
}

func (v nullValue) Kind() Kind {
	return NullKind
}

func NewEnumValue(s string) Value {
	return &enumValue{
		nameComponent: nameComponent(s),
	}
}

func (v enumValue) Value() interface{} {
	return v.Name()
}

func (v enumValue) Kind() Kind {
	return EnumKind
}

func NewObjectField(name string, value Value) ObjectField {
	return &objectField{
		nameComponent: nameComponent(name),
		valueComponent: valueComponent{value: value},
	}
}

func NewObjectValue() ObjectValue {
	return &objectValue{}
}

func (o objectValue) Kind() Kind {
	return ObjectKind
}

func (o *objectValue) Fields() chan ObjectField {
	return o.fields.Iterator()
}

func (o *objectValue) AddFields(f ...ObjectField) {
	o.fields.Add(f...)
}

// This doesn't really make sense, but... hmm, revisit
func (o objectValue) Value() interface{} {
	return nil
}
