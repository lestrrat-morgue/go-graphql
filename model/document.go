package model

func NewDocument() Document {
	return &document{}
}

func (doc *document) LookupQuery(name string) (OperationDefinition, bool) {
	doc.qmu.Lock()
	defer doc.qmu.Unlock()

	if doc.queries == nil {
		return nil, false
	}
	def, ok := doc.queries[name]
	return def, ok
}

func (doc *document) addDefinition(def Definition) {
	switch def.(type) {
	case OperationDefinition:
		var odef = def.(OperationDefinition)
		var m map[string]OperationDefinition
		switch odef.OperationType() {
		case OperationTypeQuery:
			doc.qmu.Lock()
			defer doc.qmu.Unlock()
			if doc.queries == nil {
				doc.queries = make(map[string]OperationDefinition)
			}
			m = doc.queries
		case OperationTypeMutation:
			doc.mmu.Lock()
			defer doc.mmu.Unlock()
			if doc.mutations == nil {
				doc.mutations = make(map[string]OperationDefinition)
			}
			m = doc.mutations
		}
		m[odef.Name()] = odef
	}
}

func (doc *document) AddDefinitions(list ...Definition) {
	doc.definitions.Add(list...)
	for _, def := range list {
		doc.addDefinition(def)
	}
}

func (doc document) Definitions() chan Definition {
	return doc.definitions.Iterator()
}
