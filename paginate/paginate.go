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
	params    *Params
}
type Params struct {
	Limit       int
	CurrentPage int
	DBS         *gorm.DB
	DebugQuery  bool
	SortTypes   []string
	LogOut      *io.Writer
}

func NewPaginate(p *Params) (*Paginate, error) {
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
	return &ppp, nil

}
func (p *Paginate) MakePaginate(listResult interface{}) (error) {
	//variables
	var (
		offset = 0
	)

	//get total records in table
	ch := make(chan bool, 1)
	go func() {
		if err := p.params.DBS.Model(listResult).Count(&p.Count).Error; err != nil {
			p.Log.Printf(err.Error())
		}
		ch <- true
	}()

	//awaiting count
	<-ch

	//check correct count param.page
	p.TotalPage = int(math.Ceil(float64(p.Count) / float64(p.params.Limit)))
	if p.params.CurrentPage > p.TotalPage {
		return errors.New("wrong page, page > totalpage")
	}

	//make offset
	if p.params.CurrentPage > 0 {
		offset = (p.params.CurrentPage - 1) * p.params.Limit
	}

	//check filters sorts
	if len(p.params.SortTypes) > 0 {
		for _, x := range p.params.SortTypes {
			p.params.DBS = p.params.DBS.Order(x)
		}
	}

	//get result
	if err := p.params.DBS.Limit(p.params.Limit).Offset(offset).Find(listResult).Error; err != nil {
		return err
	}
	p.Records = listResult

	//return result
	return nil
}
func (p *Paginate) Reconfig(newconfig *Params) {
	p.params = newconfig
}
