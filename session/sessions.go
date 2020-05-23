package session

import "encoding/json"

//определение типа сессии
type Session map[string]map[string]interface{}

//создание новой сессии
func NewSession() Session {
	return make(map[string]map[string]interface{})

}

//извлечение объекта из сессии
func (s *Session) Get(section, key string) interface{} {
	if res, exists := (*s)[section]; exists {
		if v, found := res[key]; found {
			return v
		}
	}
	return nil
}

//помещение объекта в сессию
func (s *Session) Put(section, key string, value interface{}) {
	if _, exists := (*s)[section]; exists {
		if _, found := (*s)[section][key]; found == false {
			(*s)[section][key] = make(map[string]interface{})
		}
	} else {
		(*s)[section] = make(map[string]interface{})
		(*s)[section][key] = make(map[string]interface{})
	}
	(*s)[section][key] = value
}

//конвертация Session -> JSON = []byte
func (s *Session) ConvertToJSON() (error, []byte) {
	if res, err := json.Marshal(s); err == nil {
		return nil, res
	} else {
		return err, nil
	}
}

//конвертация  string(JSON) -> Session
func (s *Session) ConvertFromJSON(v []byte) (error, *Session) {
	var ns = new(Session)
	if err := json.Unmarshal(v, ns); err != nil {
		return err, nil
	}
	return nil, ns
}
