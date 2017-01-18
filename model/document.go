package model

func NewDocument() *Document {
	return &Document{}
}

func (doc *Document) AddDefinitions(defs ...Definition) {
	doc.definitions = append(doc.definitions, defs...)
}

func (doc Document) Definitions() chan Definition {
	ch := make(chan Definition, len(doc.definitions))
	for _, def := range doc.definitions {
		ch <- def
	}
	close(ch)
	return ch
}
