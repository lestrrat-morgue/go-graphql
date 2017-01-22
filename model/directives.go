package model

func NewDirective(name string) Directive {
	return &directive{
		name: name,
	}
}

func (d directive) Name() string {
	return d.name
}

func (d directive) Arguments() chan Argument {
	return d.arguments.Iterator()
}

func (d *directive) AddArguments(args ...Argument) {
	d.arguments.Add(args...)
}
