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
	errorPage = errors.New("not correct number current page")
	errorType = errors.New("error: type is not slice")
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
	var pr PaginateResult
	var tt = reflect.TypeOf(list)
	var vv = reflect.ValueOf(list)

	//check slice type
	if tt.Kind() != reflect.Slice {
		return nil, errorType
	}
	//make list interface for export
	var res = make([]interface{}, vv.Len())
	for i := 0; i < vv.Len(); i++ {
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

	// calculate totalpage
	pr.Page = page
	yes := vv.Len() % countPage
	if yes > 0 {
		pr.TotalPage = (vv.Len() / countPage) + 1
	} else {
		pr.TotalPage = vv.Len() / countPage
	}

	//check correct number page
	if page <= 0 || page > pr.TotalPage {
		return nil, errorPage
	}

	//all correct
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

	//make listLinks
	var tmp = make([]int, countLink+1)
	var arr = make([]int, pr.TotalPage)
	for i := 0; i < pr.TotalPage; i++ {
		arr[i] = i
	}

	if countLink >= pr.TotalPage {
		tmp = arr
	} else {
		var right = 0
		var left = 0
		//right
		if page+countLink >= pr.TotalPage {
			right = pr.TotalPage
		} else {
			right = page + countLink
		}
		//left
		if page-(countLink-1) <= 0 {
			left = 0
		} else {
			left = page - (countLink - 1)
		}
		tmp = arr[left:right]
	}
	pr.ListPage = tmp

	//return result
	pr.CountPage = countPage
	pr.List = res
	pr.Block = result[page-1]
	pr.TotalBlock = result
	return &pr, nil
}

//
//if countLink  <= TotalPage:
//show ALL(coountLink)
//if countLink > TotalPage:
//midcountLinks  = countLink / 2 (среднее значение)
//if (CurrentPage + midcountLinks) > TotalPage:
//rightPos = [currentPage:...]
//if (CurrentPage + midcountLinks) <= TotalPage:
//rightPos = CurrentPage + midcountLinks
//if (currentPage - midcountLinks) <= 1:
//leftPos = [:currentPage]
//if (currentPage - midcountLinks) > 1:
//leftPos =  currentPage - midcountLinks
