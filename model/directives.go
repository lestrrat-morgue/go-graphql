package model

func (d *DirectiveList) Add(directives ...*Directive) {
	*d = append(*d, directives...)
}

func (d DirectiveList) Iterator() chan *Directive {
	ch := make(chan *Directive, len(d))
	for _, dir := range d {
		ch <- dir
	}
	close(ch)
	return ch
}

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
