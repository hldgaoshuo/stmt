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

func (s *SymbolTable) Set(name string) {
	s.Store[name] = s.NumDefinitions
	s.NumDefinitions++
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
