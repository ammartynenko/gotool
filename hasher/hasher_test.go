package hasher

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	h := New()
	fmt.Println(h)
}
func TestHasher_Encode(t *testing.T) {
	h := New()
	for _, x := range []string{"t", "te", "tes", "test", "tester"} {
		hashpassword := h.Encode(x)
		t.Log("===>", hashpassword)
	}
}
func TestHasher_CheckHashPassword(t *testing.T) {
	h := New()
	for _, x := range []string{"t", "te", "tes", "test", "tester"} {
		hashpassword := h.Encode(x)
		statuscheck, err := h.CheckHashPassword(hashpassword, x)
		if err != nil {
			t.Error(err)
		} else {
			t.Log("===>", hashpassword, statuscheck)
		}
	}
}
