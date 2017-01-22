package model

func NewArgument(name string, value Value) Argument {
	return &argument{
		nameComponent:  nameComponent(name),
		valueComponent: valueComponent{value: value},
	}
}
