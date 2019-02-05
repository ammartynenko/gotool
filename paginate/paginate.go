//---------------------------------------------------------------------------
//  paginate for ORM  Gorm
//---------------------------------------------------------------------------
package paginate

import (
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"io"
	"math"
	"errors"
	"strconv"
)

const prefix = "[gotool][paginate]"

type Paginate struct {
	TotalPage int
	Count     int
	Records   interface{}
	Page      int
	Log       *log.Logger
	Params    *Params
	Help      *HTMLPaginate
}
type HTMLPaginate struct {
	Totalpage   string
	Currentpage string
	Predpage    string
	Nextpage    string
	List        []string //количество элементов в пагинации
}
type Params struct {
	Limit      int
	DBS        *gorm.DB
	DebugQuery bool
	SortTypes  []string
	LogOut     *io.Writer
}

func NewPaginate(p *Params) (*Paginate) {
	var ppp Paginate
	//logger
	if p.LogOut == nil {
		ppp.Log = log.New(os.Stdout, prefix, log.Lshortfile|log.Ldate|log.Ltime)
	} else {
		ppp.Log = log.New(*p.LogOut, prefix, log.Lshortfile|log.Ldate|log.Ltime)
	}

	//config
	ppp.Params = p

	//check debug
	if p.DebugQuery {
		p.DBS = p.DBS.Debug()
	}

	//add help
	ppp.Help = &HTMLPaginate{}

	//return instance
	return &ppp

}
func (p *Paginate) MakePaginate(page int, listResult interface{}) (error) {

	//variables
	var (
		offset = 0
	)
	if page == 0 {
		p.Page = 1
	} else {
		p.Page = page
	}

	//get total records in table
	ch := make(chan bool, 1)

	go func() {
		if err := p.Params.DBS.Model(listResult).Count(&p.Count).Error; err != nil {
			p.Log.Println(err)
		}
		ch <- true
	}()

	//awaiting count
	<-ch

	//check correct count param.page
	p.TotalPage = int(math.Ceil(float64(p.Count) / float64(p.Params.Limit)))

	if p.TotalPage == 0 {
		p.TotalPage = 1
	} else {
		if p.Page > p.TotalPage {
			return errors.New("wrong page, page > totalpage")
		}
	}

	//make offset
	if p.Page > 0 {
		offset = (p.Page - 1) * p.Params.Limit
	}

	//check filters sorts
	if len(p.Params.SortTypes) > 0 {
		for _, x := range p.Params.SortTypes {
			p.Params.DBS = p.Params.DBS.Order(x)
		}
	}

	//get result
	if err := p.Params.DBS.Limit(p.Params.Limit).Offset(offset).Find(listResult).Error; err != nil {
		return err
	}
	p.Records = listResult

	//htmlhelp
	p.Help = &HTMLPaginate{}

	if p.Page == 0 {
		p.Help.Currentpage = "1"
		p.Help.Predpage = "1"
	} else {
		p.Help.Currentpage = strconv.Itoa(p.Page)
	}
	if p.Page == 1 {
		p.Help.Predpage = "1"
	} else if p.Page > 1 {
		p.Help.Predpage = strconv.Itoa(p.Page - 1)
	}

	if p.Page < p.TotalPage {
		p.Help.Nextpage = strconv.Itoa(p.Page + 1)
	} else if p.Page == p.TotalPage {
		p.Help.Nextpage = strconv.Itoa(p.Page)
	}

	p.Help.Totalpage = strconv.Itoa(p.TotalPage)
	for i := 1; i <= p.TotalPage; i++ {
		p.Help.List = append(p.Help.List, strconv.Itoa(i))
	}

	//return result
	return nil
}
func (p *Paginate) Reconfig(newconfig *Params) {
	p.Params = newconfig
}
