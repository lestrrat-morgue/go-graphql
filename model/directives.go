package model

func NewDirective(name string) *Directive {
	return &Directive{
		name: name,
	}
}

func (d Directive) Name() string {
	return d.name
}

func (d Directive) Arguments() chan *Argument {
	return d.arguments.Iterator()
}

func (d *Directive) AddArguments(args ...*Argument) {
	d.arguments.Add(args...)
}
