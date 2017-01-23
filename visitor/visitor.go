package visitor

import (
	"github.com/lestrrat/go-graphql/model"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
)

// Handler is the container for all handler functions that may be
// called while visiting a graphql data structure. You may choose
// to populate only the fields that you are interested in.
type Handler struct {
	// EnterSchema is called when starting to visit a model.Schema node.
	// All of the type definitions, including the query, are visited
	// afterward
	EnterSchema func(context.Context, model.Schema) error

	// LeaveSchema is called when leaving a model.Schema node.
	LeaveSchema func(context.Context, model.Schema) error

	// EnterDocument is called when starting to visit an Document node.
	EnterDocument func(context.Context, model.Document) error

	// LeaveDocument is called when leaving a model.Document node.
	LeaveDocument func(context.Context, model.Document) error

	// EnterDefinitionList is called when starting to traverse a
	// list of `model.Definition`s.
	EnterDefinitionList func(context.Context) error

	// LeaveDefinitionList is called when leaving a list of `model.Definition`s.
	LeaveDefinitionList func(context.Context) error

	// EnterDefinition is called when starting to visit an Definition node.
	// Note that this is called *BEFORE* determining the actual type of the
	// definition. If you only care about a specify definition type,
	// specify the handler for that specific definition type instead
	EnterDefinition func(context.Context, model.Definition) error

	// LeaveDefinition is called when leaving a model.Definition node.
	LeaveDefinition func(context.Context, model.Definition) error

	// EnterDirectiveList is called when starting to traverse a
	// list of `model.Directive`s.
	EnterDirectiveList func(context.Context) error

	// LeaveDirectiveList is called when leaving a list of `model.Directive`s.
	LeaveDirectiveList func(context.Context) error

	// EnterDirective is called when starting to visit an Directive node.
	// Arguments are not visited.
	// Does NOT respect the Pruner return value
	EnterDirective func(context.Context, model.Directive) error

	// LeaveDirective is called when leaving a model.Directive node.
	LeaveDirective func(context.Context, model.Directive) error

	// EnterOperationDefinition is called when starting to visit a model.OperationDefinition node
	// node. Selections within the definitions are followed. Variable definitions
	// are NOT followed.
	EnterOperationDefinition func(context.Context, model.OperationDefinition) error

	// LeaveOperationDefinition is called when leaving a model.OperationDefinition node.
	LeaveOperationDefinition func(context.Context, model.OperationDefinition) error

	// EnterVariableDefinition is called when starting to visit an VariableDefinition
	EnterFragmentDefinition func(context.Context, model.FragmentDefinition) error

	// LeaveVariableDefinition is called when leaving a model.VariableDefinition node.
	LeaveFragmentDefinition func(context.Context, model.FragmentDefinition) error

	// EnterObjectDefinition is called when starting to visit a model.ObjectDefinition node.
	EnterObjectDefinition func(context.Context, model.ObjectDefinition) error

	// LeaveObjectDefinition is called when leaving a model.ObjectDefinition node.
	LeaveObjectDefinition func(context.Context, model.ObjectDefinition) error

	// EnterObjectFieldDefinitionList is called when starting to traverse a
	// list of `model.ObjectFieldDefinition`s.
	EnterObjectFieldDefinitionList func(context.Context) error

	// LeaveObjectFieldDefinitionList is called when leaving a list of `model.ObjectFieldDefinition`s.
	LeaveObjectFieldDefinitionList func(context.Context) error

	// EnterObjectFieldDefinition is called when starting to visit a model.ObjectFieldDefinition node.
	EnterObjectFieldDefinition func(context.Context, model.ObjectFieldDefinition) error

	// LeaveObjectFieldDefinition is called when leaving a model.ObjectFieldDefinition node.
	LeaveObjectFieldDefinition func(context.Context, model.ObjectFieldDefinition) error

	// EnterInterfaceDefinition is called when starting to visit a model.InterfaceDefinition node.
	EnterInterfaceDefinition func(context.Context, model.InterfaceDefinition) error

	// LeaveInterfaceDefinition is called when leaving a model.InterfaceDefinition node.
	LeaveInterfaceDefinition func(context.Context, model.InterfaceDefinition) error

	// EnterInterfaceFieldDefinition is called when starting to visit a model.InterfaceFieldDefinition node.
	EnterInterfaceFieldDefinition func(context.Context, model.InterfaceFieldDefinition) error

	// LeaveInterfaceFieldDefinition is called when leaving a model.InterfaceFieldDefinition node.
	LeaveInterfaceFieldDefinition func(context.Context, model.InterfaceFieldDefinition) error

	// EnterEnumDefinition is called when starting to visit an EnumDefinition node.
	EnterEnumDefinition func(context.Context, model.EnumDefinition) error

	// LeaveEnumDefinition is called when leaving a model.EnumDefinition node.
	LeaveEnumDefinition func(context.Context, model.EnumDefinition) error

	// EnterUnionDefinition is called when starting to visit an UnionDefinition node.
	EnterUnionDefinition func(context.Context, model.UnionDefinition) error

	// LeaveUnionDefinition is called when leaving a model.UnionDefinition node.
	LeaveUnionDefinition func(context.Context, model.UnionDefinition) error

	// EnterInputDefinition is called when starting to visit a model.InputDefinition node.
	EnterInputDefinition func(context.Context, model.InputDefinition) error

	// LeaveInputDefinition is called when leaving a model.InputDefinition node.
	LeaveInputDefinition func(context.Context, model.InputDefinition) error

	// EnterInputFieldDefinitionList is called when starting to traverse a
	// list of `model.InputFieldDefinition`s.
	EnterInputFieldDefinitionList func(context.Context) error

	// LeaveInputFieldDefinitionList is called when leaving a list of `model.InputFieldDefinition`s.
	LeaveInputFieldDefinitionList func(context.Context) error

	// EnterInputFieldDefinition is called when starting to visit a model.InputFieldDefinition node.
	EnterInputFieldDefinition func(context.Context, model.InputFieldDefinition) error

	// LeaveInputFieldDefinition is called when leaving a model.InputFieldDefinition node.
	LeaveInputFieldDefinition func(context.Context, model.InputFieldDefinition) error

	// EnterSelectionList is called when starting to traverse a
	// list of `model.Selection`s.
	EnterSelectionList func(context.Context) error

	// LeaveSelectionList is called when leaving a list of `model.Selection`s.
	LeaveSelectionList func(context.Context) error

	// EnterSelection is called when starting to visit a model.Selection node.
	// Note that this is called *BEFORE* determining the actual type of the
	// selection. If you only care about a specify definition type,
	// specify the handler for that specific definition type instead
	EnterSelection func(context.Context, model.Selection) error

	// LeaveSelection is called when leaving a model.Selection node.
	LeaveSelection func(context.Context, model.Selection) error

	// EnterSelectionField is called when starting to visit a model.SelectionField node.
	EnterSelectionField func(context.Context, model.SelectionField) error

	// LeaveSelectionField is called when leaving a model.SelectionField node.
	LeaveSelectionField func(context.Context, model.SelectionField) error

	// EnterFragmentSpread is called when starting to visit a model.FragmentSpread node.
	EnterFragmentSpread func(context.Context, model.FragmentSpread) error

	// LeaveFragmentSpread is called when leaving a model.FragmentSpread node.
	LeaveFragmentSpread func(context.Context, model.FragmentSpread) error

	// EnterInlineFragment is called when starting to visit a model.InlineFragment node.
	EnterInlineFragment func(context.Context, model.InlineFragment) error

	// LeaveInlineFragment is called when leaving a model.InlineFragment node.
	LeaveInlineFragment func(context.Context, model.InlineFragment) error
}

// Pruner is the interface for errors that tell the visitor to prune
// the child nodes or not.
//
// When an EnterXXXX handler is called, the value returned from the
// Prune() method will be respected when deciding to visit the child nodes
// or not.
//
// The corresponding LeaveXXXX handler will be called regardless of the
// return value from Prune()
type Pruner interface {
	Prune() bool
}

func isPruneError(err error) (Pruner, bool) {
	if p, ok := err.(Pruner); ok {
		return p, true
	}
	return nil, false
}

// Visit starts visiting the given node structure, and calls the appropriate
// handlers that are registered in the `h` argument
func Visit(ctx context.Context, h *Handler, v interface{}) error {
	switch v.(type) {
	case model.Document:
		return visitDocument(ctx, h, v.(model.Document))
	case model.Schema:
		return visitSchema(ctx, h, v.(model.Schema))
	}
	return errors.Errorf(`invalid input type for visit: %T`, v)
}

func visitSchema(ctx context.Context, h *Handler, v model.Schema) error {
	var prune bool
	if hfunc := h.EnterSchema; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit document (enter)`)
			}
		}
	}

	if !prune {
		// This is hackisch, but we need to combine the query and the
		// definition list in one iterator in order to properly
		// handle the various cases
		typch := v.Types()
		ch := make(chan model.Definition, len(typch)+1)
		for e := range typch {
			ch <- e
		}
		ch <- v.Query()
		close(ch)

		if err := visitDefinitionList(ctx, h, ch); err != nil {
			return errors.Wrap(err, `failed to visit schema component list`)
		}
	}

	if hfunc := h.LeaveSchema; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit document (leave)`)
		}
	}
	return nil
}

func visitDocument(ctx context.Context, h *Handler, v model.Document) error {
	var prune bool
	if hfunc := h.EnterDocument; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit document (enter)`)
			}
		}
	}

	if !prune {
		if err := visitDefinitionList(ctx, h, v.Definitions()); err != nil {
			return errors.Wrap(err, `failed to visit definition list`)
		}
	}

	if hfunc := h.LeaveDocument; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit document (leave)`)
		}
	}
	return nil
}

func visitDefinitionList(ctx context.Context, h *Handler, ch chan model.Definition) error {
	if len(ch) == 0 {
		return nil
	}
	if hfunc := h.EnterDefinitionList; hfunc != nil {
		if err := hfunc(ctx); err != nil {
			return errors.Wrap(err, `failed to handle definition list (enter)`)
		}
	}
	for def := range ch {
		if err := visitDefinition(ctx, h, def); err != nil {
			return errors.Wrap(err, `failed to visit document definition`)
		}
	}
	if hfunc := h.LeaveDefinitionList; hfunc != nil {
		if err := hfunc(ctx); err != nil {
			return errors.Wrap(err, `failed to handle definition list (leave)`)
		}
	}
	return nil
}

func visitDefinition(ctx context.Context, h *Handler, v model.Definition) error {
	var prune bool
	if hfunc := h.EnterDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit definition (enter)`)
			}
		}
	}

	if !prune {
		switch v.(type) {
		case model.OperationDefinition:
			if err := visitOperationDefinition(ctx, h, v.(model.OperationDefinition)); err != nil {
				return errors.Wrap(err, `failed to visit operation definition`)
			}
		case model.FragmentDefinition:
			if err := visitFragmentDefinition(ctx, h, v.(model.FragmentDefinition)); err != nil {
				return errors.Wrap(err, `failed to visit fragment definition`)
			}
		case model.ObjectDefinition:
			if err := visitObjectDefinition(ctx, h, v.(model.ObjectDefinition)); err != nil {
				return errors.Wrap(err, `failed to visit object type definition`)
			}
		case model.InterfaceDefinition:
			if err := visitInterfaceDefinition(ctx, h, v.(model.InterfaceDefinition)); err != nil {
				return errors.Wrap(err, `failed to visit object type definition`)
			}
		case model.EnumDefinition:
			if err := visitEnumDefinition(ctx, h, v.(model.EnumDefinition)); err != nil {
				return errors.Wrap(err, `failed to visit enum definition`)
			}
		case model.UnionDefinition:
			if err := visitUnionDefinition(ctx, h, v.(model.UnionDefinition)); err != nil {
				return errors.Wrap(err, `failed to visit union definition`)
			}
		case model.InputDefinition:
			if err := visitInputDefinition(ctx, h, v.(model.InputDefinition)); err != nil {
				return errors.Wrap(err, `failed to visit input definition`)
			}
		default:
			return errors.Errorf(`unknown definition %T`, v)
		}
	}
	if hfunc := h.LeaveDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit definition (leave)`)
		}
	}
	return nil
}

func visitOperationDefinition(ctx context.Context, h *Handler, v model.OperationDefinition) error {
	var prune bool
	if hfunc := h.EnterOperationDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit operation definition (enter)`)
			}
		}
	}

	if !prune {
		if err := visitSelectionList(ctx, h, v.Selections()); err != nil {
			return errors.Wrap(err, `failed to visit selection list`)
		}
	}

	if hfunc := h.LeaveOperationDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit operation definition (leave)`)
		}
	}
	return nil
}

func visitSelectionList(ctx context.Context, h *Handler, ch chan model.Selection) error {
	if len(ch) == 0 {
		return nil
	}
	if hfunc := h.EnterSelectionList; hfunc != nil {
		if err := hfunc(ctx); err != nil {
			return errors.Wrap(err, `failed to handle selection list (enter)`)
		}
	}
	for sel := range ch {
		if err := visitSelection(ctx, h, sel); err != nil {
			return errors.Wrap(err, `failed to visit selection`)
		}
	}
	if hfunc := h.LeaveSelectionList; hfunc != nil {
		if err := hfunc(ctx); err != nil {
			return errors.Wrap(err, `failed to handle selection list (leave)`)
		}
	}
	return nil
}

func visitSelection(ctx context.Context, h *Handler, v model.Selection) error {
	var prune bool
	if hfunc := h.EnterSelection; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit selection (enter)`)
			}
		}
	}

	if !prune {
		switch v.(type) {
		case model.SelectionField:
			if err := visitSelectionField(ctx, h, v.(model.SelectionField)); err != nil {
				return errors.Wrap(err, `failed to visit selection field`)
			}
		case model.FragmentSpread:
			if err := visitFragmentSpread(ctx, h, v.(model.FragmentSpread)); err != nil {
				return errors.Wrap(err, `failed to visit fragment spread`)
			}
		case model.InlineFragment:
			if err := visitInlineFragment(ctx, h, v.(model.InlineFragment)); err != nil {
				return errors.Wrap(err, `failed to visit inline fragment`)
			}
		default:
			return errors.Errorf(`invalid selection type %T`, v)
		}
	}

	if hfunc := h.LeaveSelection; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit selection (leave)`)
		}
	}
	return nil
}

func visitSelectionField(ctx context.Context, h *Handler, v model.SelectionField) error {
	var prune bool
	if hfunc := h.EnterSelectionField; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit selection field (enter)`)
			}
		}
	}

	if !prune {
		if err := visitDirectiveList(ctx, h, v.Directives()); err != nil {
			return errors.Wrap(err, `failed to visit directive list`)
		}

		if err := visitSelectionList(ctx, h, v.Selections()); err != nil {
			return errors.Wrap(err, `failed to visit selection list`)
		}
	}

	if hfunc := h.LeaveSelectionField; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit selection field (leave)`)
		}
	}
	return nil
}

func visitDirectiveList(ctx context.Context, h *Handler, ch chan model.Directive) error {
	if len(ch) == 0 {
		return nil
	}

	if hfunc := h.EnterDirectiveList; hfunc != nil {
		if err := hfunc(ctx); err != nil {
			return errors.Wrap(err, `failed to visit directive list (enter)`)
		}
	}

	for dir := range ch {
		if err := visitDirective(ctx, h, dir); err != nil {
			return errors.Wrap(err, `failed to visit directive`)
		}
	}

	if hfunc := h.LeaveDirectiveList; hfunc != nil {
		if err := hfunc(ctx); err != nil {
			return errors.Wrap(err, `failed to visit directive list (leave)`)
		}
	}
	return nil
}

func visitDirective(ctx context.Context, h *Handler, v model.Directive) error {
	if hfunc := h.EnterDirective; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit directive (enter)`)
		}
	}

	if hfunc := h.LeaveDirective; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit directive (leave)`)
		}
	}
	return nil
}

func visitFragmentSpread(ctx context.Context, h *Handler, v model.FragmentSpread) error {
	if hfunc := h.EnterFragmentSpread; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit fragment spread (enter)`)
		}
	}
	if hfunc := h.LeaveFragmentSpread; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit fragment spread (leave)`)
		}
	}
	return nil
}

func visitInlineFragment(ctx context.Context, h *Handler, v model.InlineFragment) error {
	var prune bool
	if hfunc := h.EnterInlineFragment; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit inline fragment (enter)`)
			}
		}
	}

	if !prune {
		if err := visitDirectiveList(ctx, h, v.Directives()); err != nil {
			return errors.Wrap(err, `failed to visit directive list`)
		}

		if err := visitSelectionList(ctx, h, v.Selections()); err != nil {
			return errors.Wrap(err, `failed to visit selection list`)
		}
	}

	if hfunc := h.LeaveInlineFragment; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit inline fragment (leave)`)
		}
	}
	return nil
}

func visitFragmentDefinition(ctx context.Context, h *Handler, v model.FragmentDefinition) error {
	var prune bool
	if hfunc := h.EnterFragmentDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit fragment definition (enter)`)
			}
		}
	}

	if !prune {
		if err := visitDirectiveList(ctx, h, v.Directives()); err != nil {
			return errors.Wrap(err, `failed to visit directive list`)
		}

		if err := visitSelectionList(ctx, h, v.Selections()); err != nil {
			return errors.Wrap(err, `failed to visit selection list`)
		}
	}

	if hfunc := h.LeaveFragmentDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit fragment definition (leave)`)
		}
	}
	return nil
}

/*
func visitObjectDefinitionList(ctx context.Context, h *Handler, ch chan model.ObjectDefinition) error {
	if len(ch) == 0 {
		return nil
	}

	if hfunc := h.EnterObjectDefinitionList; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
				return errors.Wrap(err, `failed to visit object definition list (enter)`)
		}
	}
		for field := range ch {
			if err := visitObjectFieldDefinition(ctx, h, field); err != nil {
				return errors.Wrap(err, `failed to visit object field definition`)
			}
		}
	if hfunc := h.LeaveObjectDefinitionList; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit object definition list (leave)`)
		}
	}
	return nil
}
*/

func visitObjectDefinition(ctx context.Context, h *Handler, v model.ObjectDefinition) error {
	var prune bool
	if hfunc := h.EnterObjectDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit object definition (enter)`)
			}
		}
	}

	if !prune {
		if err := visitObjectFieldDefinitionList(ctx, h, v.Fields()); err != nil {
			return errors.Wrap(err, `failed to visit object definition list`)
		}
	}

	if hfunc := h.LeaveObjectDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit object definition (leave)`)
		}
	}
	return nil
}

func visitObjectFieldDefinitionList(ctx context.Context, h *Handler, ch chan model.ObjectFieldDefinition) error {
	if len(ch) == 0 {
		return nil
	}

	if hfunc := h.EnterObjectFieldDefinitionList; hfunc != nil {
		if err := hfunc(ctx); err != nil {
			return errors.Wrap(err, `failed to visit object field definition list (enter)`)
		}
	}

	for field := range ch {
		if err := visitObjectFieldDefinition(ctx, h, field); err != nil {
			return errors.Wrap(err, `failed to visit object field definition`)
		}
	}

	if hfunc := h.LeaveObjectFieldDefinitionList; hfunc != nil {
		if err := hfunc(ctx); err != nil {
			return errors.Wrap(err, `failed to visit object field definition list (leave)`)
		}
	}
	return nil
}

func visitObjectFieldDefinition(ctx context.Context, h *Handler, v model.ObjectFieldDefinition) error {
	if hfunc := h.EnterObjectFieldDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit object field definition (enter)`)
		}
	}

	if hfunc := h.LeaveObjectFieldDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit object field definition (leave)`)
		}
	}

	return nil
}

func visitInterfaceDefinition(ctx context.Context, h *Handler, v model.InterfaceDefinition) error {
	var prune bool
	if hfunc := h.EnterInterfaceDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit interface definition (enter)`)
			}
		}
	}

	if !prune {
		for field := range v.Fields() {
			if err := visitInterfaceFieldDefinition(ctx, h, field); err != nil {
				return errors.Wrap(err, `failed to visit interface field definition`)
			}
		}
	}

	if hfunc := h.LeaveInterfaceDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit interface definition (leave)`)
		}
	}
	return nil
}

func visitInterfaceFieldDefinition(ctx context.Context, h *Handler, v model.InterfaceFieldDefinition) error {
	if hfunc := h.EnterInterfaceFieldDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit interface field definition (enter)`)
		}
	}

	if hfunc := h.LeaveInterfaceFieldDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit interface field definition (leave)`)
		}
	}

	return nil
}

func visitEnumDefinition(ctx context.Context, h *Handler, v model.EnumDefinition) error {
	if hfunc := h.EnterEnumDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit enum definition (enter)`)
		}
	}

	if hfunc := h.LeaveEnumDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit enum definition (leave)`)
		}
	}
	return nil
}

func visitUnionDefinition(ctx context.Context, h *Handler, v model.UnionDefinition) error {
	if hfunc := h.EnterUnionDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit union definition (enter)`)
		}
	}

	if hfunc := h.LeaveUnionDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit union definition (leave)`)
		}
	}
	return nil
}

func visitInputDefinition(ctx context.Context, h *Handler, v model.InputDefinition) error {
	var prune bool
	if hfunc := h.EnterInputDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit input definition (enter)`)
			}
		}
	}

	if !prune {
		if err := visitInputFieldDefinitionList(ctx, h, v.Fields()); err != nil {
			return errors.Wrap(err, `failed to visit input field definition list`)
		}
	}

	if hfunc := h.LeaveInputDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit input definition (leave)`)
		}
	}
	return nil
}

func visitInputFieldDefinitionList(ctx context.Context, h *Handler, ch chan model.InputFieldDefinition) error {
	var prune bool
	if hfunc := h.EnterInputFieldDefinitionList; hfunc != nil {
		if err := hfunc(ctx); err != nil {
			if perr, ok := isPruneError(err); ok {
				prune = perr.Prune()
			} else {
				return errors.Wrap(err, `failed to visit input field definition list (enter)`)
			}
		}
	}

	if !prune {
		for e := range ch {
			if err := visitInputFieldDefinition(ctx, h, e); err != nil {
				return errors.Wrap(err, `failed to visit input field definition`)
			}
		}
	}

	if hfunc := h.EnterInputFieldDefinitionList; hfunc != nil {
		if err := hfunc(ctx); err != nil {
			return errors.Wrap(err, `failed to visit input field definition list (enter)`)
		}
	}
	return nil
}

func visitInputFieldDefinition(ctx context.Context, h *Handler, v model.InputFieldDefinition) error {
	if hfunc := h.EnterInputFieldDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit input field definition (enter)`)
		}
	}

	if hfunc := h.LeaveInputFieldDefinition; hfunc != nil {
		if err := hfunc(ctx, v); err != nil {
			return errors.Wrap(err, `failed to visit input field definition (leave)`)
		}
	}

	return nil
}
