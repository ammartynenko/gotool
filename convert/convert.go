//---------------------------------------------------------------------------
//  конвертер - пакет для помоши в конвертации разных величин и приведения их
//	к приемлемому типу для ситуации
//---------------------------------------------------------------------------
package convert

import (
	"fmt"
	d "github.com/mozillazg/go-unidecode"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Convert struct {
	logger   *log.Logger
	value    interface{}
	result   interface{}
	stockFu  map[string]func()
	Validate []int
	Replacer []int
	InValid  []int
}

const (
	LAYOUT_HTML_UTC        = "2006-01-02 15:04:05 -0700 MST"
	LAYOUT_HTML_DATA       = "2006-01-02"
	LAYOUT_HTML_RFC3339    = "2006-01-02T15:04"
	LAYOUT_TIME_LAYOUT     = "15:04"
	LAYOUT_DATE_LAYOUT     = "2006-01-02"
	LAYOUT_DATETIME_LAYOUT = "2006-01-02 15:04"
)

var (
	acceptTypes []interface{} = []interface{}{
		"", 0, int64(0),
	}
	letterBytes               = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	defConverter              = "[goconverter] %s\n"
	ErrorValueNotValidConvert = "error convert value"
)

func NewConverter(lg *log.Logger) *Convert {
	f := &Convert{
		stockFu: make(map[string]func()),
	}
	if lg == nil {
		f.logger = log.New(os.Stdout, defConverter, log.Ltime|log.Ldate)
	} else {
		f.logger = lg
	}
	f.stockFu["string"] = f.stringToInt
	f.stockFu["string"] = f.stringToInt64
	//	hex 65-122 A-z (допустимые )
	//	hex 48-57 0-9 ( допустимые )
	// 	hex 20 (зменяемые)
	//	hex 123-126, 91-96, 58-64, 33-47 punctuations ( запретные )
	f.Validate = f.convert(65, 122)
	f.InValid = f.convert(123, 126)
	f.InValid = append(f.InValid, f.convert(91, 96)...)
	f.InValid = append(f.InValid, f.convert(58, 64)...)
	f.InValid = append(f.InValid, f.convert(33, 47)...)
	f.InValid = append(f.InValid, f.convert(33, 47)...)
	f.Replacer = f.convert(32, 32)
	return f
}

//конвертация строки в целоцисленное значение 32 разрядное
func (m *Convert) StrToInt() *Convert {
	if f, exists := m.stockFu["string"]; exists {
		f()
	}
	return m
}

//конвертация строки в целоцисленное значение 64 разрядное
func (m *Convert) StrToInt64() *Convert {
	if f, exists := m.stockFu["string"]; exists {
		f()
	}
	return m
}

//конвертация строки в целоцисленное значение 64 разрядное
func (m *Convert) stringToInt64() {
	m.stringToInt()
	if m.result != nil {
		m.result = int64(m.result.(int))
	} else {
		m.result = nil
	}
}

//конвертация строки в целоцисленное значение 32 разрядное
func (m *Convert) stringToInt() {
	if r, err := strconv.Atoi(m.value.(string)); err != nil {
		m.logger.Printf(defConverter, err.Error())
		m.result = nil
	} else {
		m.result = r
	}
}

//возвращает результат последней конвертации
func (m *Convert) Result() interface{} {
	return m.result
}

//  инциализация вводным значением
func (m *Convert) Value(value interface{}) *Convert {
	if m.checkValue(value) {
		m.value = value
		return m
	}
	return nil
}

//  проверка типа поступившего значения на возможность конвертации
func (m *Convert) checkValue(value interface{}) bool {
	tValue := reflect.TypeOf(value)
	for _, x := range acceptTypes {
		if tValue == reflect.TypeOf(x) {
			return true
		}
	}
	m.logger.Printf(defConverter, ErrorValueNotValidConvert)
	return false
}

//конвертация плавающего значения в строку
func (m *Convert) FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 2, 64)
}

//конвертация целочисленного 64 разрядного значения в строку
func (m *Convert) Int64ToString(input_num int64) string {
	return strconv.FormatInt(input_num, 10)

}
//конвертация строки в uint
func (m *Convert) DirectStringToUint(v string) uint {
	if res, err := strconv.ParseUint(v, 10, 32); err != nil {
		m.logger.Println(err)
		return 0
	}  else {
		return uint(res)
	}
}


//конвертация строки в целочисленное 64 разрядное число
func (m *Convert) DirectStringtoInt64(v string) int64 {
	if res, err := strconv.Atoi(v); err != nil {
		m.logger.Printf(defConverter, err.Error())
		return 0
	} else {
		return int64(res)
	}
}

//конертация строки в список целочисленных значений
func (m *Convert) DirectStringtoIntSlice(a []string) []int {
	var result []int
	if len(a) > 0 {
		for _, x := range a {
			if res, err := strconv.Atoi(x); err != nil {
				m.logger.Printf(defConverter, err.Error())
				continue
			} else {
				result = append(result, res)
			}
		}
	}
	return result
}

//конвертация строки в список целочисленных 64 разрядных чисел
func (m *Convert) DirectStringtoInt64Slice(a []string) []int64 {
	var result []int64
	if len(a) > 0 {
		for _, x := range a {
			if res, err := strconv.Atoi(x); err != nil {
				m.logger.Printf(defConverter, err.Error())
				continue
			} else {
				result = append(result, int64(res))
			}
		}
	}
	return result
}

//конвертация строки в булево значение
func (m *Convert) DirectStringFormtoBool(v string) bool {
	if v == "" {
		return false
	}
	return true
}

//конвертация строки в целочисленное число
func (m *Convert) DirectStringtoInt(v string) int {
	if len(v) > 0 {
		if res, err := strconv.Atoi(v); err != nil {
			m.logger.Printf(defConverter, err.Error())
			return 0
		} else {
			return res
		}
	}
	return 0

}

//конвертация строки в плавающее 64 разрядное число
func (m *Convert) DirectStringtoFloat64(v string) float64 {
	if res, err := strconv.ParseFloat(v, 10); err != nil {
		m.logger.Printf(defConverter, err.Error())
		return 0
	} else {
		return res
	}
}

//=========================================
// DATA CONVERT
//=========================================
//из числового значения времени переводим в строковое представление раздельное даты и времени
type StockTime struct {
	Date, Time string
}

func (s StockTime) String() string {
	return fmt.Sprintf("%s %s", s.Date, s.Time)
}
//конвертация из time.time в строкове значение для представления в форме
//из базы данных(приложения) [time.Time] -> html (input.data, input.time) [string]
func (ms *Convert) TimeToForm(v time.Time) StockTime {
	y, m, d := v.Date()
	var (
		Less   = "0%d"
		More   = "%d"
		month  string
		day    string
		hour   string
		minute string
		st     StockTime
	)
	//data:month
	if m < 10 {
		month = fmt.Sprintf(Less, m)
	} else {
		month = fmt.Sprintf(More, m)
	}

	//data:day
	if d < 10 {
		day = fmt.Sprintf(Less, d)
	} else {
		day = fmt.Sprintf(More, d)
	}
	st.Date = fmt.Sprintf("%v-%v-%v", y, month, day)

	//time
	if v.Hour() < 10 {
		hour = fmt.Sprintf(Less, v.Hour())
	} else {
		hour = fmt.Sprintf(More, v.Hour())
	}
	if v.Minute() < 10 {
		minute = fmt.Sprintf(Less, v.Minute())
	} else {
		minute = fmt.Sprintf(More, v.Minute())
	}
	st.Time = strings.Join([]string{hour, minute}, ":")

	return st
}

//конвертация из time.unix в строковое значение для представления в форме
func (ms *Convert) UnixToForm(unixTime int64) StockTime {
	ut := time.Unix(unixTime, 0)
	y, m, d := ut.Date()
	var (
		Less   = "0%d"
		More   = "%d"
		month  string
		day    string
		hour   string
		minute string
		st     StockTime
	)
	//data:month
	if m < 10 {
		month = fmt.Sprintf(Less, m)
	} else {
		month = fmt.Sprintf(More, m)
	}

	//data:day
	if d < 10 {
		day = fmt.Sprintf(Less, d)
	} else {
		day = fmt.Sprintf(More, d)
	}
	st.Date = fmt.Sprintf("%v-%v-%v", y, month, day)

	//time
	if ut.Hour() < 10 {
		hour = fmt.Sprintf(Less, ut.Hour())
	} else {
		hour = fmt.Sprintf(More, ut.Hour())
	}
	if ut.Minute() < 10 {
		minute = fmt.Sprintf(Less, ut.Minute())
	} else {
		minute = fmt.Sprintf(More, ut.Minute())
	}
	st.Time = strings.Join([]string{hour, minute}, ":")

	return st
}

//конвертация из формы строковых значений в тип time.unix
func (ms *Convert) FormToUnix(s StockTime, correctOffset time.Duration) (int64, error) {
	nt, err := time.Parse(LAYOUT_DATETIME_LAYOUT, s.String())
	if err != nil {
		ms.logger.Printf("%v", err)
		return 0, err
	}
	nt = nt.Add(correctOffset)
	return nt.Unix(), nil
}

//конвертация из формы строковых значений в тип времени
// [string=input.data][string=input.time] -> time.Time ["2006-10-10 15:02"]
func (ms *Convert) FormToTime(s StockTime, correctOffset time.Duration) (time.Time, error) {
	nt, err := time.Parse(LAYOUT_DATETIME_LAYOUT, s.String())
	if err != nil {
		ms.logger.Printf("%v", err)
		return time.Time{}, err
	}
	nt = nt.Add(correctOffset)
	return nt, nil
}

//конвертация unix int64 в time.time
func (ms *Convert) UnixToTime(unixTime int64, correctOffset time.Duration) time.Time {
	return time.Unix(unixTime, 0).Add(correctOffset)
}

//конвертация database=TimeStamp (LAYOUT_HTML_UTC) в HTML-datetime-local(string)=LAYOUT_HTML_RFC3339
func (m *Convert) StringUTCtoHTML3339string(v string) string {
	return m.StringUTCtoDate(v).Format(LAYOUT_HTML_RFC3339)
}

//конвертация UTC в time. (html DATA из формы конвертируется этой функцией)
//(FROM HTML)html.input(type=datetime-local) -> time.Time
func (m *Convert) StringUTCtoDate(o string) time.Time {
	//layout := "2006-01-02 15:04:05 -0700 MST"
	t, err := time.Parse(LAYOUT_HTML_UTC, o)
	if err != nil {
		m.logger.Fatal(err)
	}
	return t
}

//конвертация HTML.DATA в time.time
//(FROM HTML) html.input(type=data) -> time.Time
func (m *Convert) StringDATAtoTime(o string) time.Time {
	t, err := time.Parse(LAYOUT_HTML_DATA, o)
	if err != nil {
		m.logger.Fatal(err)
	}
	return t
}

//конвертация HTML.DATA в time.time
//(FROM HTML:: RFC3389) html.input(type=datatime-local) -> time.Time
func (m *Convert) StringDATA3389toTime(o string) time.Time {
	t, err := time.Parse(LAYOUT_HTML_RFC3339, o)
	if err != nil {
		m.logger.Fatal(err)
	}
	return t
}

//конвертация Time.time в LAYOUT_HTML_DATA
//(TO HTML) time.time -> html.input(type=data)
func (m *Convert) TimeToDATA(o time.Time) string {
	return o.Format(LAYOUT_HTML_DATA)
}

//конвертация Time.time в HTML_RFC3389
//(TO HTML) time.time -> html.input(type=datatime-local)
func (m *Convert) TimeToDATARFC3389(o time.Time) string {
	return o.Format(LAYOUT_HTML_RFC3339)
}

func (m *Convert) TimeToTimeDATA(o time.Time) time.Time {
	tt := m.StringDATAtoTime(o.Format(LAYOUT_HTML_DATA))
	return tt
}

//конвертация Time.time в LAYOUT_HTML_DATA
//(TO HTML) time.time -> html.input(type=datetime-local)
func (m *Convert) TimeToDATA_UTC(o time.Time) string {
	return o.Format(LAYOUT_HTML_UTC)
}

// конертация HTML даты в Unix формат
func (m *Convert) ConvertHTMLDatetoUnix(date string) (int64, error) {
	if len(date) > 0 {
		result, err := time.Parse(LAYOUT_HTML_DATA, date)
		if err == nil {
			return result.Unix(), err
		} else {
			return 0, err
		}
	}
	return 0, nil
}

//конвертация UNIX временного предсталения в строку
func (m *Convert) ConvertUnixTimeToString(unixtime int64) string {
	return time.Unix(unixtime, 0).String()
}

//convert timeUnix->HTML5Datatime_local(string)
func (m *Convert) TimeUnixToDataLocal(unixtime int64) string {
	tmp_result := time.Unix(unixtime, 0).Format(time.RFC3339)
	g := strings.Join(strings.SplitAfterN(tmp_result, ":", 3)[:2], "")
	return g[:len(g)-1]
}

//convert HTML5Datatime_local(string)->TimeUnix
func (m *Convert) DataLocalToTimeUnix(datatimeLocal string) int64 {
	r, _ := time.Parse(time.RFC3339, datatimeLocal+":00Z")
	return r.Unix()
}

//convert HTML5Data->UnixTime
func (m *Convert) HTML5DataToUnix(s string) int64 {
	l := "2006-01-02"
	r, _ := time.Parse(l, s)
	return r.Unix()
}

//convert HTML5Data->time.Time
func (m *Convert) HTML5DataToTime(s string) time.Time {
	l := "2006-01-02T15:04"
	r, _ := time.Parse(l, s)
	return r
}

//UnixTime->HTML5Data
func (m *Convert) UnixtimetoHTML5Date(unixtime int64) string {
	return time.Unix(unixtime, 0).Format("2006-01-02")
}

//рандомный генератор строк переменной длины
func (m *Convert) RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func (s *Convert) convert(start, end int) []int {
	stack := []int{}
	for ; start <= end; start++ {
		stack = append(stack, start)
	}
	return stack
}
func (s *Convert) correct(str string) string {

	var result []string
	for _, x := range strings.Split(strings.TrimSpace(str), " ") {
		if x != "" {
			result = append(result, x)
		}
	}
	return strings.Join(result, " ")
}
func (s *Convert) preCorrect(str string) string {
	str = s.correct(str)
	var tmp []string
	for _, sym := range str {
		switch {
		case s.InSlice(s.InValid, int(sym)):
			continue
		case s.InSlice(s.Validate, int(sym)):
			tmp = append(tmp, string(sym))
		case s.InSlice(s.Replacer, int(sym)):
			tmp = append(tmp, " ")
		default:
			tmp = append(tmp, string(sym))
		}
	}
	return s.correct(strings.Join(tmp, ""))
}

//конвертация-транслитерация имени файла
func (s *Convert) TransliterCyrFilename(filename string) string {
	var extension = filepath.Ext(filename)
	var name = filename[0 : len(filename)-len(extension)]

	name = s.preCorrect(name)
	var result []string
	for _, sym := range d.Unidecode(name) {
		switch {
		case s.InSlice(s.InValid, int(sym)):
			continue
		case s.InSlice(s.Validate, int(sym)):
			result = append(result, string(sym))
		case s.InSlice(s.Replacer, int(sym)):
			result = append(result, "-")
		default:
			result = append(result, string(sym))
		}
	}
	return strings.Join([]string{strings.Join(result, ""), extension}, "")
}

//конвертация-транслитерация в кириллическое представление параметра функции
func (s *Convert) TransliterCyr(str string) string {
	str = s.preCorrect(str)
	var result []string
	for _, sym := range d.Unidecode(str) {
		switch {
		case s.InSlice(s.InValid, int(sym)):
			continue
		case s.InSlice(s.Validate, int(sym)):
			result = append(result, string(sym))
		case s.InSlice(s.Replacer, int(sym)):
			result = append(result, "-")
		default:
			result = append(result, string(sym))
		}
	}
	return strings.Join(result, "")
}

//проверка вхождения целочисленного элемента в список элементов
func (s *Convert) InSlice(str []int, target int) bool {
	for x := 0; x < len(str); x++ {
		if str[x] == target {
			return true
		}
	}
	return false
}
func (s *Convert) ShowAscii() {
	var i int
	for i = 0; i < 255; i++ {
		fmt.Printf("Dec: %3d Sym: %3c Hex: %3x\n", i, i, i)
	}
}

//---------------------------------------------------------------------------
//  convert string to TIme.Time
//---------------------------------------------------------------------------
func (s *Convert) StringToTime(year, mont, day, hour, minute, second int) *time.Time {
	layout2 := "2006-01-02 15:04:05"
	var res []string
	for _, x := range []int{year, mont, day, hour, minute, second} {
		if x < 10 {
			res = append(res, fmt.Sprintf("0%d", x))
		} else {
			res = append(res, fmt.Sprintf("%d", x))
		}
	}

	t, err := time.Parse(layout2, fmt.Sprintf("%s-%s-%s %s:%s:%s", res[0], res[1], res[2], res[3], res[4], res[5]))
	if err != nil {
		s.logger.Printf(err.Error())
		return nil
	}
	return &t
}

//---------------------------------------------------------------------------
//  check type elemnent [bool/int] and return bool result
//--------------------------------------------------------
func (s *Convert) ThisInt(v interface{}) bool {
	switch v.(type) {
	case int, int8, int16, int32, int64:
		return true
	}
	return false
}
func (s *Convert) ThisString(v interface{}) bool {
	switch v.(type) {
	case string:
		return true
	}
	return false
}
func (s *Convert) ThisFloat(v interface{}) bool {
	switch v.(type) {
	case float32, float64:
		return true
	}
	return false
}
func (s *Convert) ThisComplex(v interface{}) bool {
	switch v.(type) {
	case complex64, complex128:
		return true
	}
	return false
}
func (s *Convert) ThisSliceINT(v interface{}) bool {
	switch v.(type) {
	case []int, []int64, []int32, []int16, []int8:
		return true
	}
	return false
}
func (s *Convert) ThisSliceString(v interface{}) bool {
	switch v.(type) {
	case []string:
		return true
	}
	return false
}
func (s *Convert) ThisStringToInt(v string) (*int, bool) {
	as, err := strconv.Atoi(v)
	if err != nil {
		return nil, false
	}
	return &as, true
}

//---------------------------------------------------------------------------
//  converter human view size bytes [bytes,kbytes,mbytes,gigabytes,terabytes,petabytes
//---------------------------------------------------------------------------
type HumanSizer struct {
	key     rune
	pattern string
	total   float64
	valid   bool
}

func NewHumaneSizer() *HumanSizer {
	return &HumanSizer{}
}
func (h HumanSizer) Total() float64 {
	return h.total
}
func (h HumanSizer) String() string {
	if h.valid {
		return fmt.Sprintf(h.pattern, h.total)
	}
	return fmt.Sprint(h.pattern)
}

func (h *HumanSizer) HumanConvert(v float64, s rune) float64 {
	switch s {
	case 'b', 'B':
		h.total = v
		h.key = 'b'
		h.pattern = "%.3f byte(s)"
		h.valid = true
	case 'k', 'K':
		h.total = v / math.Pow(1024, 1)
		h.key = 'k'
		h.pattern = "%.3f KBytes"
		h.valid = true
	case 'm', 'M':
		h.total = v / math.Pow(1024, 2)
		h.key = 'm'
		h.pattern = "%.3f MBytes"
		h.valid = true
	case 'g', 'G':
		h.total = v / math.Pow(1024, 3)
		h.key = 'g'
		h.pattern = "%.3f GBytes"
		h.valid = true
	case 't', 'T':
		h.total = v / math.Pow(1024, 4)
		h.key = 't'
		h.pattern = "%.3f TBytes"
		h.valid = true
	case 'p', 'P':
		h.total = v / math.Pow(1024, 5)
		h.key = 'p'
		h.pattern = "%.3f PBytes"
		h.valid = true
	default:
		log.Fatal("wrong type size")
		h.pattern = "[WRONG TYPE FOR VIEW]"
		h.valid = false
	}
	return h.total
}
