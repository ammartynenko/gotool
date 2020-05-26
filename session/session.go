//
// другая модель сессии, с разграничением текстового блока key[string]:value[string], и бинарного блока, представленного
// key[string]:value[interface] под эгидой одной структуры, даст возможность дампить текстового блока в базе данных
// или экспортировать в формате JSON etc...
// бинарный блок даст передавать в едином контексте инстанса Сессии нужный функционал в шаблоны
// НО: кто мешает использовать нужный функционал непосредственно в обработчиках или передавать непосредственно
// в шаблоны
// в шаблоны можно пулять только функции, все остальное  идет как контейнер
// предполагается, что инстанс будет браться из пула sync.pool, поэтому мьютексы не использую
//
package session

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

const (
	logPREFIX  = "* [session.v2] "
	logBITMASK = log.LstdFlags | log.Lshortfile
)

var (
	errorNOTFOUND  = errors.New("not found key")
	errorNOTFOUNDS = errors.New("not found section")
)

//текстовой блок
type TEXT map[string]map[string]interface{}

//"бинарный" блок
type DATA map[string]map[string]interface{}

//инстанс сессии
type Session struct {
	TEXT  //interface -> для большей гибкости при передачи данных, т.к. json корректно кодирует некотоорые вложенные структуры данных
	DATA  //"бинарный" контейнер, используется только в рамках контекста, не сохраняется
	FLASH //флэш
	log   *log.Logger
}

func New() *Session {
	return &Session{
		TEXT:  make(TEXT),
		DATA:  make(DATA),
		FLASH: newFlash(),
		log:   log.New(os.Stdout, logPREFIX, logBITMASK),
	}
}

//размещения данных в ~ТЕКСТЕ~
func (t *TEXT) set(s, k string, v interface{}) {
	if _, exists := (*t)[s]; exists == false {
		(*t)[s] = make(map[string]interface{})
	}
	(*t)[s][k] = v
}

//извлечение данных из ~ТЕКСТА~
func (t *TEXT) get(s, k string) (interface{}, error) {
	if _, exists := (*t)[s]; exists {
		if r, found := (*t)[s][k]; found {
			return r, nil
		} else {
			return nil, errorNOTFOUND
		}
	}
	return nil, errorNOTFOUNDS
}

//конвертация TEXT  -> JSON = []byte
func (t *TEXT) toJSON() ([]byte, error) {
	if res, err := json.Marshal(t); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

//конвертация  string(JSON) -> Session.text
func (t *TEXT) fromJSON(v []byte) (TEXT, error) {
	var ns = make(map[string]map[string]interface{})
	if err := json.Unmarshal(v, &ns); err != nil {
		return nil, err
	}
	return ns, nil
}

//размещения данных в ~DATA~
func (t *DATA) set(s, k string, v interface{}) {
	if _, exists := (*t)[s]; exists == false {
		(*t)[s] = make(map[string]interface{})
	}
	(*t)[s][k] = v
}

//извлечение данных из ~DATA~
func (t *DATA) get(s, k string) (interface{}, error) {
	if _, exists := (*t)[s]; exists {
		if r, found := (*t)[s][k]; found {
			return r, nil
		} else {
			return nil, errorNOTFOUND
		}
	}
	return nil, errorNOTFOUNDS
}

//размещение в ТЕКСТовом блоке
func (s *Session) SetTEXT(section, key string, value interface{}) {
	s.TEXT.set(section, key, value)
}

//размещение в DATA блоке
func (s *Session) SetDATA(section, key string, value interface{}) {
	s.DATA.set(section, key, value)
}

//извлечение из ТЕКСТового блока
func (s *Session) GetTEXT(section, key string) (interface{}, error) {
	return s.TEXT.get(section, key)
}

//извлечение из DATA блока
func (s *Session) GetDATA(section, key string) (interface{}, error) {
	return s.DATA.get(section, key)
}

//извлечение из ТЕКСТового блока
func (s *Session) TextToJSON() ([]byte, error) {
	return s.TEXT.toJSON()
}

//из строки в ТЕКСТ
func (s *Session) TextFROMJSON(v []byte) (TEXT, error) {
	return s.TEXT.fromJSON(v)
}

//создание Session.Data
func (s *Session) NewDATA() DATA {
	return make(map[string]map[string]interface{})
}

//создание Session.Text
func (s *Session) NewTEXT() TEXT {
	return make(map[string]map[string]interface{})

}

//обновление TEXT
func (s *Session) UpdateTEXT(v TEXT) {
	s.TEXT = v
}

//флэш
type FLASH map[string][]interface{}

func newFlash() FLASH {
	return FLASH{}
}
func (f *FLASH) Get(key string) []interface{} {
	return (*f)[key]
}
func (f *FLASH) Set(key string, v interface{}) {
	(*f)[key] = append((*f)[key], v)
}
func (f *FLASH) Empty(key string) bool {
	return len((*f)[key]) > 0
}
