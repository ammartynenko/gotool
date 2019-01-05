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
	Params    *Params
}
type Params struct {
	Limit       int
	CurrentPage int
	DBS         *gorm.DB
	DebugQuery  bool
	SortTypes   []string
	LogOut      *io.Writer
}

func NewPaginate(p *Params) (*Paginate) {
	var ppp Paginate
	//logger
	if p.LogOut == nil {
		ppp.Log = log.New(os.Stdout, prefix, log.Lshortfile|log.Ldate|log.Ltime)
	} else {
		ppp.Log = log.New(*p.LogOut, prefix, log.Lshortfile|log.Ldate|log.Ltime)
	}

	//check debug
	if p.DebugQuery {
		p.DBS = p.DBS.Debug()
	}

	//return instance
	return &ppp

}
func (p *Paginate) MakePaginate(listResult interface{}) (error) {
	//variables
	var (
		offset = 0
	)

	//get total records in table
	ch := make(chan bool, 1)
	go func() {
		if err := p.Params.DBS.Model(listResult).Count(&p.Count).Error; err != nil {
			p.Log.Printf(err.Error())
		}
		ch <- true
	}()

	//awaiting count
	<-ch

	//check correct count param.page
	p.TotalPage = int(math.Ceil(float64(p.Count) / float64(p.Params.Limit)))
	if p.Params.CurrentPage > p.TotalPage {
		return errors.New("wrong page, page > totalpage")
	}

	//make offset
	if p.Params.CurrentPage > 0 {
		offset = (p.Params.CurrentPage - 1) * p.Params.Limit
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

	//return result
	return nil
}
func (p *Paginate) Reconfig(newconfig *Params) {
	p.Params = newconfig
}
