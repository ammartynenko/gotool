package session

//определение типа сессии
type Session map[string]map[string]interface{}

//создание новой сессии
func NewSession() Session {
	return make(map[string]map[string]interface{})
}

//извлечение объекта из сессии
func (s Session) Get(section, key string) interface{} {
	_, exists := s[section]
	if exists {
		v, found := s[section][key]
		if found {
			return v
		}
		return nil
	}
	return nil
}

//помещение объекта в сессию
func (s Session) Put(section, key string, value interface{}) {
	_, exists := s[section]
	if exists {
		_, found := s[section][key]
		if !found {
			s[section][key] = make(map[string]interface{})
		}
	} else {
		s[section] = make(map[string]interface{})
		s[section][key] = make(map[string]interface{})
	}
	s[section][key] = value
}
