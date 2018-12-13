//---------------------------------------------------------------------------
//  контейнер для сохранения состояние чего-нибудь между запросами
//---------------------------------------------------------------------------

package databox

import (
	"sync"
)

type StatelessData struct {
	sync.RWMutex
	Data map[string]*Box
}
type Box struct {
	Key  string
	Data map[string]interface{}
}

func NewDataBox() *StatelessData {
	return &StatelessData{
		Data: make(map[string]*Box),
	}
}
func (s *StatelessData) NewBox(key string) *Box {
	return &Box{
		Key:  key,
		Data: make(map[string]interface{}),
	}
}
func (s *StatelessData) Save(key string, value *Box) {
	s.Lock()
	defer s.Unlock()
	s.Data[key] = value
}
func (s *StatelessData) Get(key string) *map[string]interface{} {
	s.Lock()
	defer s.Unlock()
	value, exists := s.Data[key]
	if !exists {
		return nil
	}
	return &value.Data
}
func (s *Box) Reset() {
	s.Data = make(map[string]interface{})
	return
}
