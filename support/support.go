package support

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"time"
)

type Support struct{}

func NewSupport() *Support {
	return new(Support)
}

//создание кукиса
func (s *Support) NewCook(cookName string, salt string, r *http.Request) http.Cookie {
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

//генерация hash-sha из байтовой последовательности
func (s *Support) HashSHA(b []byte) string {
	h := sha1.New()
	h.Write(b)
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}
