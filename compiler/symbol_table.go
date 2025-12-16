package compiler

var Global *SymbolTable

const (
	GlobalScope string = "GLOBAL"
	UpScope     string = "UP"
	LocalScope  string = "LOCAL"
)

type LocalInfo struct {
	Name  string
	Index uint64
}

func NewLocalInfo(name string, index uint64) *LocalInfo {
	return &LocalInfo{
		Name:  name,
		Index: index,
	}
}

type UpInfo struct {
	LocalIndex uint64
	IsLocal    bool
}

func NewUpInfo(localIndex uint64, isLocal bool) *UpInfo {
	return &UpInfo{
		LocalIndex: localIndex,
		IsLocal:    isLocal,
	}
}

type SymbolTable struct {
	Outer       *SymbolTable
	LocalValues map[string]*LocalInfo
	UpValues    []*UpInfo
}

func NewSymbolTable(outer *SymbolTable) *SymbolTable {
	inner := &SymbolTable{
		Outer:       outer,
		LocalValues: map[string]*LocalInfo{},
		UpValues:    []*UpInfo{},
	}
	if outer == nil {
		Global = inner
	}
	return inner
}

func (s *SymbolTable) DefineGlobal(name string) error {
	if _, ex := s.LocalValues[name]; ex {
		return ErrVariableAlreadyDefined
	}
	index := uint64(len(s.LocalValues))
	localInfo := NewLocalInfo(name, index)
	s.LocalValues[name] = localInfo
	return nil
}

func (s *SymbolTable) Define(name string) (uint64, string, error) {
	if s.Outer == nil {
		localInfo := s.LocalValues[name]
		return localInfo.Index, GlobalScope, nil
	}
	if _, ex := s.LocalValues[name]; ex {
		return 0, "", ErrVariableAlreadyDefined
	}
	index := uint64(len(s.LocalValues))
	localInfo := NewLocalInfo(name, index)
	s.LocalValues[name] = localInfo
	return index, LocalScope, nil
}

func (s *SymbolTable) Get(name string) (uint64, string, bool) {
	if localInfo, ex := s.LocalValues[name]; ex {
		if s.Outer == nil {
			return localInfo.Index, GlobalScope, true
		}
		return localInfo.Index, LocalScope, true
	}
	if s.Outer == nil {
		return 0, "", false
	}

	symbolIndex, symbolScope, ex := s.Outer.Get(name)
	if !ex {
		return 0, "", false
	}
	switch symbolScope {
	case GlobalScope:
		return symbolIndex, GlobalScope, true
	case LocalScope:
		upIndex := s.UpValuesLen()
		s.UpValuesAdd(symbolIndex, true)
		return upIndex, UpScope, true
	case UpScope:
		s.UpValuesAdd(symbolIndex, false)
		return symbolIndex, UpScope, true
	default:
		return 0, "", false
	}
}

func (s *SymbolTable) UpValuesLen() uint64 {
	return uint64(len(s.UpValues))
}

func (s *SymbolTable) UpValuesAdd(symbolIndex uint64, isLocal bool) {
	upInfo := s.UpValuesIndex(symbolIndex, isLocal)
	if upInfo == nil {
		upInfo = NewUpInfo(symbolIndex, isLocal)
	}
	s.UpValues = append(s.UpValues, upInfo)
}

func (s *SymbolTable) UpValuesIndex(symbolIndex uint64, isLocal bool) *UpInfo {
	for _, upInfo := range s.UpValues {
		if upInfo.LocalIndex == symbolIndex && upInfo.IsLocal == isLocal {
			return upInfo
		}
	}
	return nil
}
