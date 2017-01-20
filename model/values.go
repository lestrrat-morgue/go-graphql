package model

import (
	"strconv"

	"github.com/pkg/errors"
)

func NewVariable(s string) *Variable {
	return &Variable{
		nameComponent: nameComponent(s),
	}
}

func (v Variable) Value() interface{} {
	return v.Name()
}

func NewIntValue(s string) (*IntValue, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse int`)
	}
	return &IntValue{
		value: int(v),
	}, nil
}

func (v IntValue) Value() interface{} {
	return v.value
}

func NewFloatValue(s string) (*FloatValue, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse float`)
	}
	return &FloatValue{
		value: v,
	}, nil
}

func (v FloatValue) Value() interface{} {
	return v.value
}

func NewStringValue(s string) *StringValue {
	return &StringValue{
		value: s,
	}
}

func (v StringValue) Value() interface{} {
	return v.value
}

func NewBoolValue(s string) (*BoolValue, error) {
	v, err := strconv.ParseBool(s)
	if err != nil {
		return nil, errors.Wrap(err, `failed to parse bool`)
	}
	return &BoolValue{
		value: v,
	}, nil
}

func (v BoolValue) Value() interface{} {
	return v.value
}

func (v NullValue) Value() interface{} {
	return nil
}

func NewEnumValue(s string) *EnumValue {
	return &EnumValue{
		nameComponent: nameComponent(s),
	}
}

func (v EnumValue) Value() interface{} {
	return v.Name()
}

func (l *ObjectFieldList) Add(v ...*ObjectField) {
	*l = append(*l, v...)
}

func (l ObjectFieldList) Iterator() chan *ObjectField {
	ch := make(chan *ObjectField, len(l))
	for _, f := range l {
		ch <- f
	}
	close(ch)
	return ch
}

func NewObjectField(name string, value Value) *ObjectField {
	return &ObjectField{
		nameComponent: nameComponent(name),
		value: value,
	}
}

func (f ObjectField) Value() Value {
	return f.value
}

func NewObjectValue() *ObjectValue {
	return &ObjectValue{}
}

func (o *ObjectValue) Fields() chan *ObjectField {
	return o.fields.Iterator()
}

func (o *ObjectValue) AddFields(f ...*ObjectField) {
	o.fields.Add(f...)
}

// This doesn't really make sense, but... hmm, revisit
func (o ObjectValue) Value() interface{} {
	return nil
}
