package support

import (
	"time"
	"crypto/md5"
	"net/http"
	"fmt"
)

type Support struct {}

func NewSupport() *Support {
	return new(Support)
}
//создание кукиса
func (s *Support) NewCook(cookName string, salt string, r *http.Request) (http.Cookie) {
	cook := http.Cookie{}
	cook.Name = cookName
	cook.Value = s.CookGenerate(salt)
	cook.Expires = time.Now().Add(time.Duration(86000*30) * time.Minute)
	cook.Path = "/"
	return cook
}
//генерация нового значения для кукиса
func (s *Support) CookGenerate(salt string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(time.Now().String()+salt)))
}

