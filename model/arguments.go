package model

func NewArgument(name string, value Value) *Argument {
	return &Argument{
		nameComponent:  nameComponent(name),
		valueComponent: valueComponent{value: value},
	}
}
