//---------------------------------------------------------------------------
//  GORENDER - рендеринг шаблонов с поддержкой "горячей" отладки, с возможностью
//  добавления собственного функционала, доступного в шаблонах
//---------------------------------------------------------------------------

package render

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"
)

const (
	PREFIXLOGGER              = "[gorender] "
	ERROR_HTTPMETHODNOTACCEPT = "[GOTOOL-RENDER] http method not allowed "
	ERROR_READTEMPLATES       = "[GOTOOL-RENDER] %s"
	ERROR_WRITETEMPLATES      = "[GOTOOL-RENDER] %s"
	ERROR_READ_TXTFILE        = "[GOTOOL-RENDER] %s"
	ERROR_WRONGTEMPLATES      = "[GOTOOL-RENDER] wrong output templates %T %v\n"
	ERROR_EXECUTETEMPLATE     = "[GOTOOL-RENDER] wrong execute template %v\n"
	ERROR_WRONGIOWRITTER      = "[GOTOOL-RENDER] FATAL - WRONG io.WRITTER\n"
	ERROR_JSON                = "[GOTOOL-RENDER] ERROR JSON %s\n"
	//---------------------------------------------------------------------------
	//  CONST:HTTP-MEDIATYPES
	//---------------------------------------------------------------------------
	ApplicationJSON                  = "application/json"
	ApplicationJSONCharsetUTF8       = ApplicationJSON + "; " + CharsetUTF8
	ApplicationJavaScript            = "application/javascript"
	ApplicationJavaScriptCharsetUTF8 = ApplicationJavaScript + "; " + CharsetUTF8
	ApplicationXML                   = "application/xml"
	ApplicationXMLCharsetUTF8        = ApplicationXML + "; " + CharsetUTF8
	ApplicationForm                  = "application/x-www-form-urlencoded"
	ApplicationProtobuf              = "application/protobuf"
	ApplicationMsgpack               = "application/msgpack"
	TextHTML                         = "text/html"
	TextHTMLCharsetUTF8              = TextHTML + "; " + CharsetUTF8
	TextPlain                        = "text/plain"
	TextPlainCharsetUTF8             = TextPlain + "; " + CharsetUTF8
	MultipartForm                    = "multipart/form-data"
	//---------------------------------------------------------------------------
	//  CONST: HTTP-CHARSET
	//---------------------------------------------------------------------------
	CharsetUTF8 = "charset=utf-8"
	//---------------------------------------------------------------------------
	//  CONST:  HTTP-HEADERS
	//---------------------------------------------------------------------------
	AcceptEncoding     = "Accept-Encoding"
	Authorization      = "Authorization"
	ContentDisposition = "Content-Disposition"
	ContentEncoding    = "Content-Encoding"
	ContentLength      = "Content-Length"
	ContentType        = "Content-Type"
	Location           = "Location"
	Upgrade            = "Upgrade"
	Vary               = "Vary"
	WWWAuthenticate    = "WWW-Authenticate"
	XForwardedFor      = "X-Forwarded-For"
	XRealIP            = "X-Real-IP"
)

//---------------------------------------------------------------------------
//  список дефолтных функций, входящих в список инстанс рендера, доступных
//  в шаблонах при обработке
//---------------------------------------------------------------------------

var (
	defaultFilters = map[string]interface{}{
		"random":              randomGenerator,
		"count":               strings.Count,
		"split":               strings.Split,
		"title":               strings.Title,
		"lower":               strings.ToLower,
		"totitle":             strings.ToTitle,
		"makemap":             makeMap,
		"in":                  mapIn,
		"andlist":             andList,
		"upper":               strings.ToUpper,
		"concat":              concat,
		"unixtime":            unixtimeNormal,
		"unixtimeformat":      unixtimeNormalFormatData,
		"unixtodata":          unixtimeNormalFormatData,
		"datehtmltounix":      hTML5DataToUnix,
		"timeUnixToDataLocal": timeUnixToDataLocal,
		"dataLocalToTimeUnix": dataLocalToTimeUnix,
		"yesno":               yesNo,
		"html2": func(value string) template.HTML {
			return template.HTML(fmt.Sprint(value))
		},
		"type":        typeIs,
		"jsonconvert": jSONconvert,
	}
)

//---------------------------------------------------------------------------
//  определение типа рендера
//---------------------------------------------------------------------------
type (
	Render struct {
		sync.RWMutex
		Temp            *template.Template
		Filters         template.FuncMap
		Debug           bool
		Path            string
		logger          *log.Logger
		Lg              io.Writer
		DebugFatal      bool //вываливать в log.Fatal при ошибка рендера шаблонов/парсинга директории с шаблонами, по умолчанию false
		logwriterEnable bool
	}
)

//добавочный для гибкости по логированию
func NewRenderL(path string, debug bool, logger io.Writer, debugFatal bool) *Render {
	sf := &Render{}
	defer sf.catcherPanic()
	sf.Filters = template.FuncMap{}
	sf.AddFilters(defaultFilters)
	sf.Path = path
	sf.Debug = debug
	sf.DebugFatal = debugFatal
	if _, valid := logger.(io.Writer); valid {
		sf.Lg = logger
		sf.logwriterEnable = true
	} else {
		sf.logger = log.New(os.Stdout, PREFIXLOGGER, log.Ltime|log.Ldate|log.Lshortfile)
		log.Printf(ERROR_WRONGIOWRITTER)
		sf.logwriterEnable = false
	}
	sf.ReloadTemplate()
	return sf
}

//создание нового инстанса // сохраняю в целях совместимости
func NewRender(path string, debug bool, logger *log.Logger, debugFatal bool) *Render {
	sf := &Render{}
	defer sf.catcherPanic()
	sf.Filters = template.FuncMap{}
	sf.AddFilters(defaultFilters)
	sf.Path = path
	sf.Debug = debug
	sf.DebugFatal = debugFatal
	sf.logwriterEnable = false
	if logger != nil {
		sf.logger = logger
	} else {
		sf.logger = log.New(os.Stdout, PREFIXLOGGER, log.Ltime|log.Ldate|log.Lshortfile)
	}
	sf.ReloadTemplate()
	return sf
}

//перезагрузка дерева шаблонов
func (s *Render) ReloadTemplate() {
	defer s.catcherPanic()
	if s.Debug || s.Temp == nil {
		s.Temp = template.Must(template.New("indexstock").Funcs(s.Filters).ParseGlob(s.Path))
	}
}

//перегружает отдельный блок/шаблон для обновления данных
func (s *Render) ExecuteTemplate(name string, data interface{}, w http.ResponseWriter) {
	if err := s.Temp.ExecuteTemplate(w, name, data); err != nil {
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_EXECUTETEMPLATE, err.Error())))
		} else {
			if s.logger != nil {
				s.logger.Printf(fmt.Sprintf(ERROR_EXECUTETEMPLATE, err.Error()))
			}
		}
	}
}

//показ указанного шаблона, с указанием data-контейнера, и интерфейса вывода
func (s *Render) Render(name string, data interface{}, w interface{}) (err error) {
	defer s.catcherPanic()
	if s.Debug || s.Temp == nil {
		s.ReloadTemplate()
	}
	buf := new(bytes.Buffer)
	if err = s.Temp.ExecuteTemplate(buf, name, data); err != nil {
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_READTEMPLATES, err.Error())))
		} else {
			s.logger.Printf(fmt.Sprintf(ERROR_READTEMPLATES, err.Error()))
		}
		if s.DebugFatal {
			if s.logwriterEnable {
				_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_READTEMPLATES, err.Error())))
			} else {
				s.logger.Printf(fmt.Sprintf(ERROR_READTEMPLATES, err.Error()))
			}
			log.Fatal(err)
		}
		return
	}
	switch w.(type) {
	case http.ResponseWriter:
		resp := w.(http.ResponseWriter)
		resp.Header().Add(ContentType, TextHTMLCharsetUTF8)
		//resp.WriteHeader(http.StatusOK)
		if _, err := resp.Write(s.HTMLTrims(buf.Bytes())); err != nil {
			if s.logwriterEnable {
				_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error())))
			} else {
				s.logger.Printf(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error()))
			}
			return err
		}
	case *os.File:
		if _, err = w.(*os.File).Write(buf.Bytes()); err != nil {
			if s.logwriterEnable {
				_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error())))
			} else {
				s.logger.Printf(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error()))
			}
			return err
		}
	default:
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_WRONGTEMPLATES, nil, nil)))
		} else {
			s.logger.Printf(fmt.Sprintf(ERROR_WRONGTEMPLATES, nil, nil))
		}
		if s.DebugFatal {
			log.Fatal(err)
		}
	}
	return
}

//показ указанного шаблона, с указанием data-контейнера, и интерфейса вывода + указание http кода
func (s *Render) RenderCode(httpCode int, name string, data interface{}, w interface{}) (err error) {
	defer s.catcherPanic()
	if s.Debug || s.Temp == nil {
		s.ReloadTemplate()
	}
	buf := new(bytes.Buffer)
	if err = s.Temp.ExecuteTemplate(buf, name, data); err != nil {
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_EXECUTETEMPLATE, err.Error())))
		} else {
			s.logger.Printf(fmt.Sprintf(ERROR_EXECUTETEMPLATE, err.Error()))
		}
		if s.DebugFatal {
			log.Fatal(err)
		}
		return err
	}
	switch w.(type) {
	case http.ResponseWriter:
		resp := w.(http.ResponseWriter)
		resp.Header().Add(ContentType, TextHTMLCharsetUTF8)
		resp.WriteHeader(httpCode)
		if _, err := resp.Write(s.HTMLTrims(buf.Bytes())); err != nil {
			if s.logwriterEnable {
				_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error())))
			} else {
				s.logger.Printf(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error()))
			}
			return err
		}
	case *os.File:
		if _, err = w.(*os.File).Write(buf.Bytes()); err != nil {
			if s.logwriterEnable {
				_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error())))
			} else {
				s.logger.Printf(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error()))
			}
			return err
		}
	default:
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_WRONGTEMPLATES, nil, nil)))
		} else {
			s.logger.Printf(fmt.Sprintf(ERROR_WRONGTEMPLATES, nil, nil))
		}
		if s.DebugFatal {
			log.Fatal(err)
		}
	}
	return
}

func (s *Render) RenderTxt(httpCode int, name string, w interface{}) (err error) {
	//read txt file
	file, err := os.Open(name)
	if err != nil {
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_READ_TXTFILE, err.Error())))
		} else {
			s.logger.Printf(ERROR_READ_TXTFILE, err.Error())
		}

		if s.DebugFatal {
			log.Fatal(err)
		}
		return err
	}
	outFile, err := ioutil.ReadAll(file)
	if err != nil {
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_READ_TXTFILE, err.Error())))
		} else {
			s.logger.Printf(ERROR_READ_TXTFILE, err.Error())
		}

		if s.DebugFatal {
			log.Fatal(err)
		}
		return err
	}
	switch w.(type) {
	case http.ResponseWriter:
		resp := w.(http.ResponseWriter)
		resp.Header().Add(ContentType, TextPlain)
		resp.WriteHeader(httpCode)
		if _, err = resp.Write(outFile); err != nil {
			if s.logwriterEnable {
				_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error())))
			} else {
				s.logger.Printf(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error()))
			}
			return err
		}
	case *os.File:
		if _, err = w.(*os.File).Write(outFile); err != nil {
			if s.logwriterEnable {
				_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error())))
			} else {
				s.logger.Printf(fmt.Sprintf(ERROR_WRITETEMPLATES, err.Error()))
			}
			return err
		}
	default:
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_WRONGTEMPLATES, nil, nil)))
		} else {
			s.logger.Printf(fmt.Sprintf(ERROR_WRONGTEMPLATES, nil, nil))
		}
		if s.DebugFatal {
			log.Fatal(err)
		}
	}
	return
}

//---------------------------------------------------------------------------
//  JSON
//---------------------------------------------------------------------------
//записывает json(byte format) в responseWriter
func (s *Render) JSONB(httpcode int, b []byte, w http.ResponseWriter) error {
	w.Header().Set(ContentType, ApplicationJavaScriptCharsetUTF8)
	w.WriteHeader(httpcode)
	if _, err := w.Write(b); err != nil {
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_JSON, err.Error())))
		} else {
			s.logger.Printf(fmt.Sprintf(ERROR_JSON, err.Error()))
		}
		return err
	}
	return nil
}

//записывает json в responseWriter
func (s *Render) JSON(code int, answer interface{}, w http.ResponseWriter) (err error) {
	b, err := json.Marshal(answer)
	if err != nil {
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_JSON, err.Error())))
		} else {
			s.logger.Printf(fmt.Sprintf(ERROR_JSON, err.Error()))
		}
		return err
	}
	return s.JSONB(code, b, w)
}

//отловка паники
func (s *Render) catcherPanic() {
	msgPanic := recover()
	if msgPanic != nil && s.logger != nil {
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf(ERROR_EXECUTETEMPLATE, msgPanic)))
		} else {
			s.logger.Printf(ERROR_EXECUTETEMPLATE, msgPanic)
		}
		if s.DebugFatal {
			s.logger.Fatal(msgPanic)
		}
	}
}

//вычищает пустые строки в шаблонах при рендеринге
func (s *Render) HTMLTrims(body []byte) []byte {
	result := []string{}
	for _, line := range strings.Split(string(body), "\n") {
		if len(line) != 0 && len(strings.TrimSpace(line)) != 0 {
			result = append(result, line)
		}
	}
	return []byte(strings.Join(result, "\n"))
}

//отображение всех функций-фильтров, доступных в шаблонах
func (s *Render) ShowFiltersFuncs(out io.Writer) {
	for name, f := range s.Filters {
		if s.logwriterEnable {
			_, _ = s.Lg.Write([]byte(fmt.Sprintf("`%s`:`%v`\n", name, f)))
		} else {
			s.logger.Printf("`%s`:`%v`\n", name, f)
		}
	}
}

//---------------------------------------------------------------------------
//  дополнительный функционал деофлтный по умолчанию реализован тут,
//  он может быть расширен, посредством добавление в карту нужных функций в шаблонах
//---------------------------------------------------------------------------
func (s *Render) AddUserFilter(name string, f interface{}) {
	s.Filters[name] = f
}
func (s *Render) AddFilters(stack map[string]interface{}) {
	for k, v := range stack {
		s.Filters[k] = v
	}
}

//возращает тип аргумента
func typeIs(value interface{}) string {
	v := reflect.ValueOf(value)
	var result string
	switch v.Kind() {
	case reflect.Bool:
		result = "bool"
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		result = "integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		result = "unsigned integer"
	case reflect.Float32, reflect.Float64:
		result = "float"
	case reflect.String:
		result = "string"
	case reflect.Slice:
		result = "slice"
	case reflect.Map:
		result = "map"
	case reflect.Chan:
		result = "chan"
	default:
		result = "undefine type"
	}
	return result
}
func mapIn(value interface{}, stock interface{}) bool {
	switch value.(type) {
	case int64:
		for _, x := range stock.([]int64) {
			if x == value.(int64) {
				return true
			}
		}
	case int:
		for _, x := range stock.([]int) {
			if x == value.(int) {
				return true
			}
		}
	case string:
		for _, x := range stock.([]string) {
			if x == value.(string) {
				return true
			}
		}

	}
	return false
}
func makeMap(value ...string) []string {
	return value
}

func andList(listValues ...interface{}) bool {
	for _, v := range listValues {
		if v == nil {
			return false
		}
	}
	return true
}
func yesNo(value bool, yes, no string) string {
	if value {
		return yes
	}
	return no
}

//---------------------------------------------------------------------------
//  TIME Functions
//---------------------------------------------------------------------------
func unixtimeNormal(unixtime int64) string {
	return time.Unix(unixtime, 0).String()
}

//UnixTime->HTML5Data
func unixtimeNormalFormatData(unixtime int64) string {
	return time.Unix(unixtime, 0).Format("2006-01-02")
}

//convert HTML5Data->UnixTime
func hTML5DataToUnix(s string) int64 {
	l := "2006-01-02"
	r, _ := time.Parse(l, s)
	return r.Unix()
}

//convert timeUnix->HTML5Datatime_local(string)
func timeUnixToDataLocal(unixtime int64) string {
	tmp_result := time.Unix(unixtime, 0).Format(time.RFC3339)
	g := strings.Join(strings.SplitAfterN(tmp_result, ":", 3)[:2], "")
	return g[:len(g)-1]
}

//convert HTML5Datatime_local(string)->TimeUnix
func dataLocalToTimeUnix(datatimeLocal string) int64 {
	r, _ := time.Parse(time.RFC3339, datatimeLocal+":00Z")
	return r.Unix()
}

//---------------------------------------------------------------------------
//  рандомный генератор для корректного обновления css,js в head
//---------------------------------------------------------------------------
func randomGenerator() int {
	return rand.Intn(1000)
}

//---------------------------------------------------------------------------
//  JSON конвертация
//---------------------------------------------------------------------------

func jSONconvert(obj interface{}) string {
	buf, err := json.Marshal(obj)
	if err != nil {
		fmt.Printf(err.Error())
		return ""
	}
	return string(buf)
}

func concat(s1, s2 string) string {
	return s2 + s1
}
