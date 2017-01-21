package model

func NewDocument() Document {
	return &document{}
}

func (doc *document) AddDefinitions(list ...Definition) {
	doc.definitions.Add(list...)
}

func (doc document) Definitions() chan Definition {
	return doc.definitions.Iterator()
}
