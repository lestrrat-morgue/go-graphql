package model

func NewArgument(name string, value Value) *Argument {
	return &Argument{
		name:  name,
		value: value,
	}
}

func (arg Argument) Name() string {
	return arg.name
}

func (arg Argument) Value() Value {
	return arg.value
}
