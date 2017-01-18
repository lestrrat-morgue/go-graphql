package model

func (args *ArgumentList) Add(list ...*Argument) {
	*args = append(*args, list...)
}

func (args ArgumentList) Iterator() chan *Argument {
	ch := make(chan *Argument, len(args))
	for _, arg := range args {
		ch <- arg
	}
	close(ch)
	return ch
}

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
