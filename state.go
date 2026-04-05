package gwiz

// State is a typed key-value bag that accumulates across wizard steps.
type State interface {
	Get(key string) (any, bool)
	GetString(key string) string
	GetBool(key string) bool
	GetInt(key string) int
	GetStringSlice(key string) []string
	Set(key string, value any)
	Keys() []string
}

// newState creates a new empty State.
func newState() State {
	return &mapState{data: make(map[string]any)}
}

type mapState struct {
	data map[string]any
}

func (s *mapState) Get(key string) (any, bool) {
	v, ok := s.data[key]
	return v, ok
}

func (s *mapState) GetString(key string) string {
	v, ok := s.data[key]
	if !ok {
		return ""
	}
	str, _ := v.(string)
	return str
}

func (s *mapState) GetBool(key string) bool {
	v, ok := s.data[key]
	if !ok {
		return false
	}
	b, _ := v.(bool)
	return b
}

func (s *mapState) GetInt(key string) int {
	v, ok := s.data[key]
	if !ok {
		return 0
	}
	i, _ := v.(int)
	return i
}

func (s *mapState) GetStringSlice(key string) []string {
	v, ok := s.data[key]
	if !ok {
		return nil
	}
	ss, _ := v.([]string)
	return ss
}

func (s *mapState) Set(key string, value any) {
	s.data[key] = value
}

func (s *mapState) Keys() []string {
	keys := make([]string, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}
