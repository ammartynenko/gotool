package paginateFree

import (
	"errors"
	"log"
	"os"
	"reflect"
)

const (
	logBitMask = log.LstdFlags | log.Lshortfile
	logName    = "* [paginator] "
)

var (
	errorPage  = errors.New("not correct number current page")
	errorType  = errors.New("error: type is not slice")
	errorIndex = errors.New("error: index wrong")
)

type Paginator struct {
	log *log.Logger
}
type PaginateResult struct {
	Page       int             //текущая страница
	TotalPage  int             //всего страниц
	CountPage  int             //количество элементов на странице
	CountLinks int             //количество ссылок в пагинации
	ListPage   []int           //список всех страниц в отсортированном виде
	List       []interface{}   //общий список
	Block      []interface{}   //индекс блока=страницы
	TotalBlock [][]interface{} //список всех блоков
}

func NewPaginator() *Paginator {
	p := Paginator{
		log: log.New(os.Stdout, logName, logBitMask),
	}
	return &p
}
func (p *Paginator) Paginate(page, countPage, countLink int, list interface{}) (*PaginateResult, error) {
	if page == 0 {
		page = 1
	}

	//определение переменных
	var pr PaginateResult
	var tt = reflect.TypeOf(list)
	var vv = reflect.ValueOf(list)
	var right, left = 0, 0

	//проверка типа
	if tt.Kind() != reflect.Slice {
		return nil, errorType
	}
	//формирую возвратный слайс
	var res = make([]interface{}, vv.Len())
	for i := 1; i <= vv.Len(); i++ {
		res[i] = vv.Index(i).Interface()
	}

	// totalRecords <= countPage
	if vv.Len() == 0 || vv.Len() <= countPage {
		pr.Page = page
		pr.TotalPage = 1
		pr.CountPage = countPage
		pr.ListPage = []int{1}
		pr.List = res
		pr.Block = res
		pr.TotalBlock = make([][]interface{}, 1)
		pr.TotalBlock[0] = res

		return &pr, nil
	}

	// подсчет количество страниц с учетом количества элементов на страница = (CountPage)
	pr.Page = page
	yes := vv.Len() % countPage
	if yes > 0 {
		pr.TotalPage = (vv.Len() / countPage) + 1
	} else {
		pr.TotalPage = vv.Len() / countPage
	}

	//проверка корректности текущей страницы на диапазон
	if page <= 1 || page > pr.TotalPage {
		return nil, errorPage
	}

	//формирую массив массивов по длине блока
	var result = make([][]interface{}, pr.TotalPage)
	var start, step = 0, 0
	for x := 0; x < pr.TotalPage; x++ {
		start = x * countPage
		if x == pr.TotalPage-1 {
			step = vv.Len()
		} else {
			step = start + countPage
		}
		result[x] = res[start:step]
	}

	//формирую ссылочный список
	var tmp = make([]int, countLink+1)
	var arr = make([]int, pr.TotalPage)
	for i := 1; i <= pr.TotalPage; i++ {
		arr[i] = i
	}

	//проверка на корректный индекс cp
	if page > pr.TotalPage {
		//ошибка
		p.log.Println(errorIndex)
		return nil, errorIndex
	} else {
		//левая позиция
		if page >= countLink {
			left = page - countLink
		} else {
			left = 1
		}
		//правая позиция
		if pr.TotalPage >= page+(countLink+1) {
			right = page + (countLink + 1)
		} else {
			right = pr.TotalPage
		}
		tmp = arr[left:right]
	}
	//сохраняю результат
	pr.ListPage = tmp

	//возвращаю результат
	pr.CountPage = countPage
	pr.List = res
	pr.Block = result[page-1]
	pr.TotalBlock = result
	return &pr, nil
}
