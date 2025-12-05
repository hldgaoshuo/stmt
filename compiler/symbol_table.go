package compiler

var Global *SymbolTable

type SymbolTable struct {
	Outer          *SymbolTable
	Store          map[string]uint64
	NumDefinitions uint64
}

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	t := &SymbolTable{
		Outer:          outer,
		Store:          make(map[string]uint64),
		NumDefinitions: 0,
	}
	if outer == nil {
		Global = t
	}
	return t
}

func (s *SymbolTable) SetGlobal(name string) error {
	if _, ok := s.Store[name]; ok {
		return ErrVariableAlreadyDefined
	}
	s.Store[name] = s.NumDefinitions
	s.NumDefinitions++
	return nil
}

func (s *SymbolTable) Set(name string) (uint64, error) {
	if s.Outer == nil {
		index, ok := s.Store[name]
		if !ok {
			return 0, ErrVariableNotDefined
		}
		return index, nil
	} else {
		if _, ok := s.Store[name]; ok {
			return 0, ErrVariableAlreadyDefined
		}
		index := s.NumDefinitions
		s.Store[name] = index
		s.NumDefinitions++
		return index, nil
	}
}

func (s *SymbolTable) Get(name string) (uint64, bool) {
	if index, ok := s.Store[name]; ok {
		return index, true
	}
	if s.Outer != nil {
		return s.Outer.Get(name)
	}
	return 0, false
}
