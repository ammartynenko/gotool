package config

import (
	"os"
	"log"
	"io/ioutil"
	"encoding/gob"
	"github.com/go-yaml/yaml"
	"io"
)

const (
	logConfigPrefix           = "[config-log] "
	logConfigFlags            = log.Ldate | log.Ltime | log.Lshortfile
	msgsuccessReadDumpConfig  = "дамп конфига успешно прочитан"
	msgsuccessWriteDumpConfig = "дамп конфига успешно записан"

	dumpfileflag = os.O_CREATE | os.O_RDWR | os.O_TRUNC
	dumpfileperm = 0666
)

//---------------------------------------------------------------------------
//  example other config
//---------------------------------------------------------------------------
type Conf struct {
	Log            *log.Logger
	PathDumpconfig string
}

func (c *Conf) ReadConfig(fileConfig string, exConf interface{}) (error) {
	//открываю файл конфигурации
	f, err := os.Open(fileConfig)
	if err != nil {
		return err
	}
	//читаю файл конфига
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	//конвертирую его в структуру
	err = yaml.Unmarshal(b, exConf)
	if err != nil {
		return err
	}
	return nil
}
func NewConf(pathDumpconfig string, logout io.Writer) *Conf {
	n := &Conf{
		PathDumpconfig: pathDumpconfig,
	}
	if logout == nil {
		n.Log = log.New(os.Stdout, logConfigPrefix, logConfigFlags)
	} else {
		n.Log = log.New(logout, logConfigPrefix, logConfigFlags)
	}
	return n
}

//запись дампа конфигурационного файла
func (c *Conf) WriteDumpConfig(config interface{}) error {
	f, err := os.OpenFile(c.PathDumpconfig, dumpfileflag, dumpfileperm)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)
	err = enc.Encode(config)
	if err != nil {
		return err
	}
	if c.Log != nil {
		c.Log.Printf(msgsuccessWriteDumpConfig)
		return err
	}
	return nil
}

//читаем дамп конфигурационного файла
func (c *Conf) ReadDumpConfig(config interface{}) (error) {
	f, err := os.Open(c.PathDumpconfig)
	if err != nil {
		return err
	}
	defer f.Close()
	dec := gob.NewDecoder(f)
	err = dec.Decode(config)
	if err != nil {
		return err
	}
	c.Log.Printf(msgsuccessReadDumpConfig)
	return nil
}
