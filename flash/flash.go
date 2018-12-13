//---------------------------------------------------------------------------
//  пакет обеспечивающий создание контейнера для передачи всплывающих сообщений
//  в контексте
//---------------------------------------------------------------------------

package flash

import (
	"sync"
	"time"
	"fmt"
	"crypto/md5"
)

type Flash struct {
	sync.RWMutex
	Key   string
	Stock map[string]*FlashMessage
	Salt  string
}
type FlashMessage struct {
	Status  string
	Message interface{}
}

//---------------------------------------------------------------------------
//  FLASH:
//---------------------------------------------------------------------------
func NewFlash(salt string) *Flash {
	n := &Flash{
		Stock: make(map[string]*FlashMessage),
		Salt:  salt,
	}
	n.Key = n.generateKey()
	return n
}
func (f *Flash) generateKey() string {
	t := time.Now()
	return fmt.Sprintf("%x", md5.Sum([]byte(t.String()+f.Salt)))
}
func (f *Flash) Set(status, section string, message interface{}) {
	nm := &FlashMessage{Status: status, Message: message}
	f.Lock()
	f.Stock[section] = nm
	f.Unlock()
}
func (f *Flash) Get(section string) (*FlashMessage) {
	f.Lock()
	defer f.Unlock()
	if result, exists := f.Stock[section]; exists {
		delete(f.Stock, section)
		return result
	}
	return nil
}
func (f *Flash) HaveMsg(section string) bool {
	f.Lock()
	defer f.Unlock()
	_, exists := f.Stock[section]
	if exists {
		return true
	}
	return false
}
