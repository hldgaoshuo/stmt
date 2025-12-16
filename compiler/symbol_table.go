package compiler

var Global *SymbolTable

const (
	GlobalScope string = "GLOBAL"
	CloserScope string = "CLOSER"
	LocalScope  string = "LOCAL"
)

type SymbolInfo struct {
	Name  string
	Index uint64
}

type SymbolTable struct {
	Outer          *SymbolTable
	Store          map[string]*SymbolInfo
	NumDefinitions uint64
}

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	inner := &SymbolTable{
		Outer:          outer,
		Store:          make(map[string]*SymbolInfo),
		NumDefinitions: 0,
	}
	if outer == nil {
		Global = inner
	}
	return inner
}

func (s *SymbolTable) DefineGlobal(name string) error {
	if _, ok := s.Store[name]; ok {
		return ErrVariableAlreadyDefined
	}
	s.Store[name] = &SymbolInfo{
		Name:  name,
		Index: s.NumDefinitions,
	}
	s.NumDefinitions++
	return nil
}

func (s *SymbolTable) Define(name string) (*SymbolInfo, string, error) {
	if s.Outer == nil {
		info, ok := s.Store[name]
		if !ok {
			return nil, "", ErrVariableNotDefined
		}
		return info, GlobalScope, nil
	}

	if _, ok := s.Store[name]; ok {
		return nil, "", ErrVariableAlreadyDefined
	}
	info := &SymbolInfo{
		Name:  name,
		Index: s.NumDefinitions,
	}
	s.Store[name] = info
	s.NumDefinitions++
	return info, LocalScope, nil
}

func (s *SymbolTable) Assign(name string) (*SymbolInfo, string, error) {
	if info, ok := s.Store[name]; ok {
		if s.Outer == nil {
			return info, GlobalScope, nil
		}
		return info, LocalScope, nil
	}

	if s.Outer == nil {
		return nil, "", ErrVariableNotDefined
	}
	info, scope, err := s.Outer.Assign(name)
	if err != nil {
		return nil, "", err
	}
	if scope == GlobalScope {
		return info, GlobalScope, nil
	}
	return info, CloserScope, nil
}

func (s *SymbolTable) Get(name string) (*SymbolInfo, string, bool) {
	if info, ok := s.Store[name]; ok {
		if s.Outer == nil {
			return info, GlobalScope, true
		}
		return info, LocalScope, true
	}

	if s.Outer == nil {
		return nil, "", false
	}
	info, scope, ok := s.Outer.Get(name)
	if !ok {
		return nil, "", false
	}
	if scope == GlobalScope {
		return info, GlobalScope, true
	}
	return info, CloserScope, true
}
