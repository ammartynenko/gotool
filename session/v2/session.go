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
package v2

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
	errorWRONGTYPE = errors.New("wrong type")
)

//текстовой блок
type TEXT map[string]map[string]interface{}

//"бинарный" блок
type DATA map[string]map[string]interface{}

//инстанс сессии
type Session struct {
	TEXT //interface -> для большей гибкости при передачи данных, т.к. json корректно кодирует некотоорые вложенные структуры данных
	DATA //"бинарный" контейнер, используется только в рамках контекста, не сохраняется
	log  *log.Logger
}

func New() *Session {
	return &Session{
		TEXT: make(TEXT),
		DATA: make(DATA),
		log:  log.New(os.Stdout, logPREFIX, logBITMASK),
	}
}

//размещения данных в ~ТЕКСТЕ~
func (t *TEXT) set(s, k string, v interface{}) {
	if _, exists := (*t)[s]; exists {
		if _, found := (*t)[s][k]; found {
		} else {
			(*t)[s] = make(map[string]interface{})
		}
	} else {
		(*t) = make(map[string]map[string]interface{})
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
func (t *TEXT) toJSON() (error, []byte) {
	if res, err := json.Marshal(t); err == nil {
		return nil, res
	} else {
		return err, nil
	}
}

//конвертация  string(JSON) -> Session.text
func (t *TEXT) fromJSON(v []byte) (error, TEXT) {
	var ns = make(map[string]map[string]interface{})
	if err := json.Unmarshal(v, &ns); err != nil {
		return err, nil
	}
	return nil, ns
}

//создание Session.Text
func (s *Session) NewTEXT() TEXT {
	ns := make(map[string]map[string]interface{})
	return ns
}
