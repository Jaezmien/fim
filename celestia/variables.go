package celestia

import (
	"git.jaezmien.com/Jaezmien/fim/luna/stack"
	"git.jaezmien.com/Jaezmien/fim/spike/variable"
)

type Variable struct {
	Name string

	*variable.DynamicVariable

	Constant bool
}

// The VariableManager will be the one handling our variables.
//
// The local stack is a stack of variables, inside another stack.
// This is because our way of managing local variables is a new
// paragraph scope, and each new stack represents one scope
// deep of statements.
type VariableManager struct {
	Globals stack.Stack[*Variable]
	Locals  stack.Stack[*stack.Stack[*Variable]]
}

func NewVariableManager() *VariableManager {
	return &VariableManager{
		Globals: *stack.New[*Variable](),
		Locals:  *stack.New[*stack.Stack[*Variable]](),
	}
}

// Create a new paragraph scope.
func (m *VariableManager) PushScope() {
	m.Locals.Push(stack.New[*Variable]())
}

// Delete current paragraph scope.
func (m *VariableManager) PopScope() {
	m.Locals.Pop()
}

// Returns the depth of paragraph scopes.
func (m *VariableManager) ScopeDepth() int {
	return m.Locals.Len()
}

func (m *VariableManager) PushVariable(variable *Variable, global bool) {
	if global {
		m.Globals.Push(variable)
	} else {
		if m.ScopeDepth() == 0 {
			panic("VariableManager@PushVariable called with no variable scopes")
		}

		current := m.Locals.Peek()
		current.Push(variable)
	}
}
func (m *VariableManager) PopVariable(global bool) *Variable {
	if global {
		return *m.Globals.Pop()
	} else {
		if m.ScopeDepth() == 0 {
			panic("VariableManager@PopVariable called with no variable scopes")
		}

		current := m.Locals.Peek()
		if current.Len() == 0  {
			panic("VariableManager@PopVarible called with an empty stack")
		}
		return *current.Pop()
	}
}
func (m *VariableManager) PopVariableAmount(global bool, amount int) []*Variable {
	variables := stack.New[*Variable]()

	if global {

		for m.Globals.Len() > 0 && amount > 0 {
			variables.Push(*m.Globals.Pop())
			amount -= 1
		}

	} else {
		if m.ScopeDepth() == 0 {
			panic("VariableManager@PopVariableAmount called with no variable scopes")
		}

		current := m.Locals.Peek()

		for current.Len() > 0 && amount > 0 {
			variables.Push(*current.Pop())
			amount -= 1
		}
	}

	return variables.Flatten()
}

func (m *VariableManager) Get(name string, local bool) *Variable {
	for idx := 0; idx < m.Globals.Len(); idx += 1 {
		variable := m.Globals.PeekAt(idx)
		if variable.Name == name {
			return variable
		}
	}

	if local && m.ScopeDepth() > 0 {
		current := m.Locals.Peek()

		for idx := 0; idx < current.Len(); idx += 1 {
			variable := current.PeekAt(idx)
			if variable.Name == name {
				return variable
			}
		}
	}

	return nil
}

func (m *VariableManager) Has(name string, local bool) bool {
	return m.Get(name, local) != nil
}
