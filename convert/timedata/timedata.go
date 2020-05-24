//
// пакет для конвертации времени в строковые шаблоны и обратно
// в целом использую для конвертации в html формах
//
package timedata

import (
	"fmt"
	"log"
	"os"
	"time"
)

const (
	logBITS            = log.LstdFlags | log.Lshortfile
	logPREFIX          = "* [timedata-DateTimeConverter] "
	lAYOUT_DATE_LAYOUT = "2006-01-02"
	lAYOUT_TIME_LAYOUT = "15:04"
)

type DateTimeConverter struct {
	log *log.Logger
}

func New() *DateTimeConverter {
	return &DateTimeConverter{
		log: log.New(os.Stdout, logPREFIX, logBITS),
	}
}

// Input: `12:22`[string] Output: ``0000-01-01 12:22:00 +0000 UTC [time.Time] :: html->(input type=time)
func (c *DateTimeConverter) StringToDate(v string) (time.Time, error) {
	return time.Parse(lAYOUT_DATE_LAYOUT, v)
}

// Input: `2020-05-21`[string]  OutPut: `2020-05-21 00:00:00 +0000 UTC` [time.Time] :: html->(input type=date)
func (c *DateTimeConverter) StringToTime(v string) (time.Time, error) {
	return time.Parse(lAYOUT_TIME_LAYOUT, v)
}

// vDate: `2020-05-21 00:00:00 +0000 UTC`  vTime: `0000-01-01 12:22:00 +0000 UTC` result: `2020-05-21 12:22:00 +0000 UTC`
func (c *DateTimeConverter) TimeDateAddTime(vData, vTime time.Time) time.Time {
	return vData.Add(time.Hour*time.Duration(vTime.Hour()) + time.Minute*time.Duration(vTime.Minute()))
}

// Input: `2020-05-21 12:22:00 +0000 UTC`[time.Time]  Output: `2020-05-21` [string]
func (c *DateTimeConverter) DisboudToDataString(v time.Time) string {
	return fmt.Sprintf("%v-%v-%v", v.Year(), v.Month(), v.Day())
}

// Input: `2020-05-21 12:22:00 +0000 UTC`[time.Time]  Output: `12:22` [string]
func (c *DateTimeConverter) DisboudToTimeString(v time.Time) string {
	return fmt.Sprintf("%v-%v", v.Hour(), v.Minute())
}
