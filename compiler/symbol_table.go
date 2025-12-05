package compiler

var Global *SymbolTable

const (
	GlobalScope string = "GLOBAL"
	LocalScope  string = "LOCAL"
)

type SymbolInfo struct {
	Name  string
	Index uint64
	Scope string
}

type SymbolTable struct {
	Outer          *SymbolTable
	Store          map[string]*SymbolInfo
	NumDefinitions uint64
}

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	t := &SymbolTable{
		Outer:          outer,
		Store:          make(map[string]*SymbolInfo),
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
	s.Store[name] = &SymbolInfo{
		Name:  name,
		Index: s.NumDefinitions,
		Scope: GlobalScope,
	}
	s.NumDefinitions++
	return nil
}

func (s *SymbolTable) Set(name string) (*SymbolInfo, error) {
	if s.Outer == nil {
		info, ok := s.Store[name]
		if !ok {
			return nil, ErrVariableNotDefined
		}
		return info, nil
	} else {
		if _, ok := s.Store[name]; ok {
			return nil, ErrVariableAlreadyDefined
		}
		info := &SymbolInfo{
			Name:  name,
			Index: s.NumDefinitions,
			Scope: LocalScope,
		}
		s.Store[name] = info
		s.NumDefinitions++
		return info, nil
	}
}

func (s *SymbolTable) Get(name string) (*SymbolInfo, bool) {
	if info, ok := s.Store[name]; ok {
		return info, true
	}
	if s.Outer != nil {
		return s.Outer.Get(name)
	}
	return nil, false
}
