package model

import (
	"strconv"

	"github.com/pkg/errors"
)

func NewVariable(s string) *Variable {
	return &Variable{
		name: s,
	}
}

func (v Variable) String() string {
	return "$" + v.name
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

func (v IntValue) String() string {
	return strconv.Itoa(v.value)
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

func (v FloatValue) String() string {
	return strconv.FormatFloat(v.value, 'g', -1, 64)
}

func NewStringValue(s string) *StringValue {
	return &StringValue{
		value: s,
	}
}

func (v StringValue) String() string {
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

func (v BoolValue) String() string {
	return strconv.FormatBool(v.value)
}

func (v NullValue) String() string {
	return "null"
}

func NewEnumValue(s string) *EnumValue {
	return &EnumValue{
		name: s,
	}
}

func (v EnumValue) String() string {
	return v.name
}
