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
)

const prefix = "[gotool][paginate]"

type Paginate struct {
	TotalPage int
	Count     int
	Records   interface{}
	Page      int
	Log       *log.Logger
}
type Params struct {
	Limit       int
	CurrentPage int
	DBS         *gorm.DB
	DebugQuery  bool
	SortTypes   []string
	LogOut      *io.Writer
}

func NewPaginate(p *Params, resultList interface{}) (*Paginate, error) {
	//variables
	var (
		ppp    Paginate
		offset = 0
	)

	//logger
	if p.LogOut == nil {
		ppp.Log = log.New(os.Stdout, prefix, log.Lshortfile|log.Ldate|log.Ltime)
	} else {
		ppp.Log = log.New(*p.LogOut, prefix, log.Lshortfile|log.Ldate|log.Ltime)
	}

	//get total records in table
	ch := make(chan bool, 1)
	go func() {
		if err := p.DBS.Model(resultList).Count(&ppp.Count).Error; err != nil {
			ppp.Log.Printf(err.Error())
		}
		ch <- true
	}()

	//awaiting count
	<-ch

	//check correct count param.page
	ppp.TotalPage = int(math.Ceil(float64(ppp.Count) / float64(p.Limit)))
	if p.CurrentPage > ppp.TotalPage {
		return nil, errors.New("wrong page, page > totalpage")
	}

	//make offset
	if p.CurrentPage > 0 {
		offset = (p.CurrentPage - 1) * p.Limit
	}

	//check filters sorts
	if len(p.SortTypes) > 0 {
		for _, x := range p.SortTypes {
			p.DBS = p.DBS.Order(x)
		}
	}

	//get result
	if err := p.DBS.Limit(p.Limit).Offset(offset).Find(resultList).Error; err != nil {
		return nil, err
	}
	ppp.Records = resultList
	
	//return result
	return &ppp, nil
}
