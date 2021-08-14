package hasher

import (
	"golang.org/x/crypto/bcrypt"
	"log"
)

type Hasher struct {
	Thashpassword string
}

func New() *Hasher {
	return &Hasher{}
}
func (h *Hasher) Encode(password string) (hashpassword string) {
	h.Thashpassword = ""
	if hs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err != nil {
		log.Println(err)
		return
	} else {
		hashpassword = string(hs)
		h.Thashpassword = hashpassword
	}
	return
}
func (h *Hasher) CheckHashPassword(hashpassword, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hashpassword), []byte(password)); err != nil {
		log.Println(err)
		return false, err
	}
	return true, nil
}
