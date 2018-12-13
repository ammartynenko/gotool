//---------------------------------------------------------------------------
//  разные общие константы/переменные
//---------------------------------------------------------------------------

package consvar

import "errors"

//---------------------------------------------------------------------------
//  CONST: константы, текстовые
//---------------------------------------------------------------------------
const (
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
	//---------------------------------------------------------------------------
	//
	//--------------------------------------------------------------------------
	PREFIXLOGGER = "[gomixer] "
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	//---------------------------------------------------------------------------
	//сообщения об ошибках в формах
	//---------------------------------------------------------------------------
	ErrorUsername = "Имя пользователя ошибочно"
	ErrorPassword = "Пароль ошибочен"
	ErrorEmail    = "Почтовый адрес ошибочен"

	//---------------------------------------------------------------------------
	//`placeholder` описания для формы
	//---------------------------------------------------------------------------
	PlaceUsername = "= имя пользователя = "
	PlacePassword = "= пароль ="
	PlaceEmail    = "= почтовый адрес ="

	//---------------------------------------------------------------------------
	//ошибки
	//---------------------------------------------------------------------------
	ParseErrorInt       = errors.New("[parseform][error] ошибка парсинга `string`->`int64`")
	PTRFormError        = errors.New("[baseform][error] Ошибка, дай мне указатель на структуру для записи")
	PTRFormErrorMethods = errors.New("[baseform][error] Ошибка, отсутствует реализация интерфейса методов для получения данных из формы")
	CSRFErrorValidate   = "CSRF не валидное значение"

	//---------------------------------------------------------------------------
	//название стилей для ошибок в формах полей
	//---------------------------------------------------------------------------
	ErrorStyleForm   = "has-error"
	SuccessStyleForm = "has-success"

	//---------------------------------------------------------------------------
	//  сообщения для ошибки в формах при валидации формы
	//---------------------------------------------------------------------------
	ErrorMsgFormString   = "- поле не может быть пустым -"
	ErrorMsgFormCheckbox = "- нажмите на чекбокс, если вы не робот -"
	ErrorMsgFormBool     = "- сделайте отметку -"
	ErrorMsgFormSelect   = "- не выбран ни один из элементов -"
)
