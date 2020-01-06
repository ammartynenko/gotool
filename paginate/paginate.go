//---------------------------------------------------------------------------
//  paginate for ORM  Gorm
//---------------------------------------------------------------------------
package paginate

import (
	"errors"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"math"
	"os"
	"strconv"
)

const prefix = "[gotool][paginate]"

type Paginate struct {
	Log    *log.Logger
	Params *Params
}
type Params struct {
	Limit      int
	DBS        *gorm.DB
	DebugQuery bool
	SortTypes  []string
	LogOut     *log.Logger
	CountLinks int
}
type HTMLPaginate struct {
	Totalpage   string
	Currentpage string
	Predpage    string
	Nextpage    string
	List        []string //количество элементов в пагинации
}

type ResultPaginate struct {
	Help       *HTMLPaginate
	TotalPage  int
	Records    interface{}
	Page       int
	Count      int
	CountLinks int
	Links      []int
	LinksStr   []string
}

func NewPaginate(p *Params) (*Paginate) {
	var ppp Paginate
	//logger
	if p.LogOut == nil {
		ppp.Log = log.New(os.Stdout, prefix, log.Lshortfile|log.Ldate|log.Ltime)
	} else {
		ppp.Log = p.LogOut
	}
	//config
	ppp.Params = p

	//return instance
	return &ppp
}
//
////=========================================
//// new refactoring make paginate result
////in args: count records, limit,
////=========================================
//type ParamsNew struct {
//	CurrentPage int
//	TotalCount  int
//	LimitCount  int
//	DBS         *sql.DB
//}
//
//func (p *Paginate) MakePaginate2(par ParamsNew) (ResultPaginate, error) {
//	//result instance
//	r := ResultPaginate{
//		Help:       &HTMLPaginate{},
//		CountLinks: p.Params.CountLinks,
//	}
//	//list total pages
//	r.TotalPage = int(math.Ceil(float64(par.TotalCount) / float64(par.LimitCount)))
//	p.Log.Printf("TotalPage: %v\n", r.TotalPage)
//	if r.TotalPage == 0 {
//		r.TotalPage = 1
//	} else {
//		if r.Page > r.TotalPage {
//			return r, errors.New("wrong page, page > totalpage")
//		}
//	}
//	//
//	end, start := 0, 0
//
//	//base loop
//	for i:=1; i <= totalPage; i ++ {
//		//first condition
//		if i == 0  {
//			start = 0
//		} else {
//			start  = end
//		}
//		//two condition
//		if i * limit > len(s) {
//			end = len(s)
//		} else {
//			end  = i * limit
//		}
//		//assign page + slice with page
//		fmt.Printf("%d : %d\n", start, end)
//		stocker[i] = s[start:end]
//	}
//
//	return r, nil
//}
//
func (p *Paginate) MakePaginate(page int, listResult interface{}) (ResultPaginate, error) {
	//result instance
	r := ResultPaginate{
		Help:       &HTMLPaginate{},
		CountLinks: p.Params.CountLinks,
	}

	//check debug
	if p.Params.DebugQuery {
		p.Params.DBS = p.Params.DBS.Debug()
	}

	//variables
	var (
		offset = 0
	)
	if page == 0 {
		r.Page = 1
	} else {
		r.Page = page
	}
	//get total records in table
	ch := make(chan bool, 1)

	go func() {
		if err := p.Params.DBS.Model(listResult).Count(&r.Count).Error; err != nil {
			p.Log.Println(err)
		}
		ch <- true
	}()

	//awaiting count
	<-ch

	//check correct count param.page
	r.TotalPage = int(math.Ceil(float64(r.Count) / float64(p.Params.Limit)))

	if r.TotalPage == 0 {
		r.TotalPage = 1
	} else {
		if r.Page > r.TotalPage {
			return r, errors.New("wrong page, page > totalpage")
		}
	}

	//make offset
	if r.Page > 0 {
		offset = (r.Page - 1) * p.Params.Limit
	}

	//check filters sorts
	if len(p.Params.SortTypes) > 0 {
		for _, x := range p.Params.SortTypes {
			p.Params.DBS = p.Params.DBS.Order(x)
		}
	}

	//get result
	if err := p.Params.DBS.Limit(p.Params.Limit).Offset(offset).Find(listResult).Error; err != nil {
		return r, err
	}
	r.Records = listResult

	if r.Page == 0 {
		r.Help.Currentpage = "1"
		r.Help.Predpage = "1"
	} else {
		r.Help.Currentpage = strconv.Itoa(r.Page)
	}
	if r.Page == 1 {
		r.Help.Predpage = "1"
	} else if r.Page > 1 {
		r.Help.Predpage = strconv.Itoa(r.Page - 1)
	}

	if r.Page < r.TotalPage {
		r.Help.Nextpage = strconv.Itoa(r.Page + 1)
	} else if r.Page == r.TotalPage {
		r.Help.Nextpage = strconv.Itoa(r.Page)
	}

	r.Help.Totalpage = strconv.Itoa(r.TotalPage)
	for i := 1; i <= r.TotalPage; i++ {
		r.Help.List = append(r.Help.List, strconv.Itoa(i))
	}

	////debug testing
	//if r.CountLinks >= r.TotalPage {
	//	for x:=1; x < r.TotalPage; x ++ {
	//		r.Links = append(r.Links, x)
	//		r.LinksStr = append(r.LinksStr, strconv.Itoa(x))
	//	}
	//}
	//if r.CountLinks < r.TotalPage {
	//	r.LinksStr = r.Help.List[r.TotalPage - p.Params.CountLinks:]
	//}

	//range available pages for view
	//если меньше
	if r.CountLinks > r.TotalPage {
		for x := 1; x <= r.TotalPage; x++ {
			r.Links = append(r.Links, x)
			r.LinksStr = append(r.LinksStr, strconv.Itoa(x))
		}
	}
	//если больше
	if r.CountLinks < r.TotalPage {
		if r.Page+r.CountLinks == r.TotalPage {
			for x := r.Page; x < r.Page+r.CountLinks; x++ {
				r.Links = append(r.Links, x)
				r.LinksStr = append(r.LinksStr, strconv.Itoa(x))
			}
		}

		if r.Page+r.CountLinks < r.TotalPage {
			for x := r.Page; x < r.Page+r.CountLinks; x++ {
				r.Links = append(r.Links, x)
				r.LinksStr = append(r.LinksStr, strconv.Itoa(x))
			}
		}

		if r.Page+r.CountLinks > r.TotalPage {
			for x := (r.TotalPage - r.CountLinks) + 1; x <= (r.TotalPage-r.CountLinks)+r.CountLinks; x++ {
				r.Links = append(r.Links, x)
				r.LinksStr = append(r.LinksStr, strconv.Itoa(x))
			}
		}
	}

	//return result
	return r, nil
}
func (p *Paginate) Reconfig(newconfig *Params) {
	p.Params = newconfig
}
