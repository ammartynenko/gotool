package paginateFree

import (
	"errors"
	"log"
	"os"
)

const (
	logBitMask = log.LstdFlags | log.Lshortfile
	logName    = "* [paginator] "
)

var (
	errorPage = errors.New("not correct number current page")
)

type Paginator struct {
	log *log.Logger
}
type PaginateResult struct {
	Page       int             //текущая страница
	TotalPage  int             //всего страниц
	CountPage  int             //количество элементов на странице
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
func (p *Paginator) Paginate(page, countPage int, list []interface{}) (*PaginateResult, error) {
	var pr PaginateResult
	// totalRecords <= countPage
	if len(list) == 0 || len(list) <= countPage {
		pr.Page = page
		pr.TotalPage = 1
		pr.CountPage = countPage
		pr.List = list
		pr.Block = list
		pr.TotalBlock = make([][]interface{}, 1)
		pr.TotalBlock[0] = list
		return &pr, nil
	}

	// calculate totalpage
	yes := len(list) % countPage
	if yes > 0 {
		pr.TotalPage = (len(list) / countPage) + 1
	} else {
		pr.TotalPage = len(list) / countPage
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
			step = len(list)
		} else {
			step = start + countPage
		}
		result[x] = list[start:step]
	}

	//return result
	pr.Page = page
	pr.CountPage = countPage
	pr.List = list
	pr.Block = result[page-1]
	pr.TotalBlock = result
	return &pr, nil
}
