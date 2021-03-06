//---------------------------------------------------------------------------
//  некоторые дефолтные компоненты для форм
//---------------------------------------------------------------------------

package forms

var (
	defforms = map[string]FieldForm{
		"email": {
			Placeholder: "=введите почтовый ящик пользователя=",
			Error:       ErrorForm{Error: "Ошибка: введите почтовый ящик пользователя", ErrorClass: "errorclass"}},
		"password": {
			Placeholder: "=введите пароль =",
			Error:       ErrorForm{Error: "Ошибка: введите пароль ", ErrorClass: "errorclass"}},
		"robot": {
			DefaultValue: "robotchecker",
			Error:        ErrorForm{ErrorClass: "errorclass", Error: "Ошибка: поставьте галочку, если вы не робот"}},
		"username": {
			Placeholder: "=введите имя пользователя=",
			Error:       ErrorForm{Error: "Ошибка: введите имя пользователя", ErrorClass: "errorclass"}},
		"title": {
			Placeholder: "=введите название=",
			Error:       ErrorForm{Error: "Ошибка: поле не может быть пустым", ErrorClass: "errorclass"}},
		"body": {
			Placeholder: "=введите данные=",
			Error:       ErrorForm{Error: "Ошибка: поле не может быть пустым", ErrorClass: "errorclass"}},
		"prebody": {
			Placeholder: "=введите данные=",
			Error:       ErrorForm{Error: "Ошибка: поле не может быть пустым", ErrorClass: "errorclass"}},
		"link": {
			Placeholder: "=введите ссылку=",
			Error:       ErrorForm{Error: "Ошибка: поле не может быть пустым", ErrorClass: "errorclass"}},
		"active": {
			Placeholder: "",
			Error:       ErrorForm{Error: "Ошибка: поставьте отметку", ErrorClass: "errorclass"}},
		//----------------------------------------------------------------------------------------------

		"name": {
			Placeholder: "=введите имя =",
			Error:       ErrorForm{Error: "Ошибка: введите название", ErrorClass: "errorclass"}},
		"family": {
			Placeholder: "=введите фамилию=",
			Error:       ErrorForm{Error: "Ошибка: введите фамилию", ErrorClass: "errorclass"}},
		"age": {
			Placeholder: "=сколько вам лет=",
			Error:       ErrorForm{Error: "Ошибка: введите количество лет пользователю", ErrorClass: "errorclass"}},
		"basketids": {
			Error: ErrorForm{Error: "Ошибка: выберете хотя бы один элемент"}, DefaultValue: []int64{1, 3, 4, 7}},
		"activeuser": {
			DefaultValue: "SOMEBITCH",
			Error:        ErrorForm{ErrorClass: "errorclass", Error: "Ошибка: поставьте галочку"}},
		"range": {
			DefaultValue: "SOMERANGE",
			Error:        ErrorForm{ErrorClass: "errorclass", Error: "Ошибка: сдвиньте ползунок до конца вправо"}},
	}
)

//formDefaultValues map[string]DefaultForm = map[string]DefaultForm{
//	"Name":      DefaultForm{Placeholder: "=введите имя пользователя=", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"Username":  DefaultForm{Placeholder: "=введите имя пользователя=", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"Password":  DefaultForm{Placeholder: "=введите пароль =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"Email":     DefaultForm{Placeholder: "=введите email =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"Port":      DefaultForm{Placeholder: "=порт сервера=", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"CatName":   DefaultForm{Placeholder: "=введите название категории=", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"Title":     DefaultForm{Placeholder: "=введите заголовок =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"MetaKeys":  DefaultForm{Placeholder: "=введите SEO слова =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"MetaDesc":  DefaultForm{Placeholder: "=введите SEO описание-сниппет =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"MetaRobot": DefaultForm{Placeholder: "=введите занчения для SEO robot=", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"Message":   DefaultForm{Placeholder: "=введите текст сообщения=", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"Body":      DefaultForm{Placeholder: "=введите тело сообщения =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"Link":      DefaultForm{Placeholder: "=введите ссылку-ключ =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"Age":       DefaultForm{Placeholder: "=введите ваш возраст=", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//
//	"UserInfoUsername": DefaultForm{Placeholder: "=введите имя пользователя=", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"UserInfoPassword": DefaultForm{Placeholder: "=введите пароль =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"UserEmail":        DefaultForm{Placeholder: "=введите email =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"CategoryName":     DefaultForm{Placeholder: "=введите название категории=", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"PostTitle":        DefaultForm{Placeholder: "=введите заголовок =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"PostBody":         DefaultForm{Placeholder: "=введите тело сообщения =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"PostSeoMetaKeys":  DefaultForm{Placeholder: "=введите SEO слова =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"PostSeoMetaDesc":  DefaultForm{Placeholder: "=введите SEO описание-сниппет =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"PostSeoMetaRobot": DefaultForm{Placeholder: "=введите занчения для SEO robot=", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"TagName":          DefaultForm{Placeholder: "=введите имя метки =", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поле не может быть пустым"},
//	"PostCategoryID":   DefaultForm{Placeholder: "", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "значение должно быть выбрано"},
//	"PostUserID":       DefaultForm{Placeholder: "", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "значение должно быть выбрано"},
//	"Robot":            DefaultForm{Placeholder: "", ErrorClassStyle: "has-error", SuccessClassStyle: "ok", ErrorMsg: "поставьте отметку что вы не робот"},
//}
