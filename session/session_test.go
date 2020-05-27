package session

import "testing"

func TestNew(t *testing.T) {
	tt := New()
	if tt.TEXT == nil || tt.DATA == nil || tt.log == nil {
		t.Error("error make new instance\n")
	}
}
func TestSession_NewDATA(t *testing.T) {
	tt := New()
	data := tt.NewDATA()
	if data == nil {
		t.Error("error make new `data`")
	}
}
func TestSession_NewTEXT(t *testing.T) {
	tt := New()
	text := tt.NewTEXT()
	if text == nil {
		t.Error("error make new `text`")
	}
}
func TestSession_GetDATA(t *testing.T) {
	tt := New()
	if _, err := tt.GetDATA("stock", "some"); err == nil {
		t.Error("error : FOUND key, byt key not make before")
	}
	tt.SetDATA("stock", "some", 123)
	if _, err := tt.GetDATA("stock", "some"); err != nil {
		t.Error("error : NOT FOUND, but stock.some -> have value")
	}
}
func TestSession_SetDATA(t *testing.T) {
	tt := New()
	tt.SetDATA("stock", "fn", 123)
	tt.SetDATA("stock", "fn2", 12312)
	tt.SetDATA("stock", "fn3", 111)
	tt.SetDATA("stock2", "fn", 123)
	t.Log(tt.DATA)
	if v, err := tt.GetDATA("stock", "fn"); err != nil {
		t.Error(err, "element before placed in stock.fnm but not found now")
	} else {
		switch rt := v.(type) {
		case int:
		default:
			t.Errorf("wrong value type, waiting `int`, getting `%T`", rt)
		}
	}
}
func TestSession_GetTEXT(t *testing.T) {
	tt := New()
	if _, err := tt.GetTEXT("stock", "simple"); err == nil {
		t.Errorf("error: waiting `NOT FOUND SECTION  AND KEY`, but result = found")
	}
	tt.SetTEXT("stock", "simple", "value")
	if _, err := tt.GetTEXT("stock", "simple"); err != nil {
		t.Errorf(err.Error(), " waiting get value, but get error")
	}
}
func TestSession_TextToJSON(t *testing.T) {
	tt := New()
	type rt struct {
		id        int64
		name      string
		radius    float64
		signature []byte
	}
	tt.SetTEXT("stock", "simple", "value")
	tt.SetTEXT("stock", "simple2", 123)
	tt.SetTEXT("stock", "simple", []string{"1", "2", "3"})
	tt.SetTEXT("stock", "simple", rt{
		id:        1,
		name:      "simpleRT",
		radius:    3.14,
		signature: []byte{1, 33, 54, 66, 123, 44, 67},
	})
	if res, err := tt.TextToJSON(); err != nil {
		t.Error(err)
	} else {
		t.Log(res)
	}
}
func TestSession_TextFROMJSON(t *testing.T) {
	tt := New()
	type rt struct {
		id        int64
		name      string
		radius    float64
		signature []byte
	}
	tt.SetTEXT("stock", "simple", "value")
	tt.SetTEXT("stock", "simple2", 123)
	tt.SetTEXT("stock", "simple", []string{"1", "2", "3"})
	tt.SetTEXT("stock", "simple", rt{
		id:        1,
		name:      "simpleRT",
		radius:    3.14,
		signature: []byte{1, 33, 54, 66, 123, 44, 67},
	})
	if res, err := tt.TextToJSON(); err != nil {
		t.Error(err)
	} else {
		if newTEXT, err := tt.TextFROMJSON(res); err != nil {
			t.Error(err)
		} else {
			t.Log(newTEXT)
		}
	}
}
func TestSession_UpdateTEXT(t *testing.T) {
	tt := New()
	tt.SetTEXT("stock", "simple", 123)
	newText := tt.NewTEXT()
	t.Log(tt.TEXT)
	newText.set("newstock", "newsimple", 456)
	tt.UpdateTEXT(newText)
	t.Log(tt.TEXT)
}

func TestFLASH_Get(t *testing.T) {
	ff := newFlash()
	ff.Set("user", "Simple test message", Success)
	ff.Set("user", "Simple test message2", Warning)

	ff.Set("Rty", "WDFDFGDFg", Success)

	t.Log(ff)
	for _, v :=  range ff.Get("user") {
		t.Log(v)
	}

	t.Log(ff)

}