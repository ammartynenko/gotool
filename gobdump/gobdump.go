//---------------------------------------------------------------------------
//  dump (read/write) for simple config struct files for example in gob format
//---------------------------------------------------------------------------

package gobdump

import (
	"encoding/gob"
	"os"
	"bytes"
)

type Configdump struct {
	b   bytes.Buffer
	enc *gob.Encoder
	dec *gob.Decoder
}

//make new instanse
func NewConfigdump() *Configdump {
	c := &Configdump{}
	c.b = bytes.Buffer{}
	c.enc = gob.NewEncoder(&c.b)
	c.dec = gob.NewDecoder(&c.b)
	return c
}

//encode interface to filename, filename all time recreate if exists, or create new
func (c *Configdump) encodefile(v interface{}, filename string) error {
	err := c.enc.Encode(v)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(c.b.Bytes())
	if err != nil {
		return err
	}
	return nil
}

//encode -> *Configdump.b (bytes.Buffer)
func (c *Configdump) encode(v interface{}) error {
	err := c.enc.Encode(v)
	if err != nil {
		return err
	}
	return nil
}

//decode -> *Configdump.b (bytes.Buffer)
func (c *Configdump) decode(v interface{}) error {
	err := c.dec.Decode(v)
	if err != nil {
		return err
	}
	return nil
}

//decode from file, file open/read and decode to interface
func (c *Configdump) decodefile(v interface{}, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	//decode filedata from new gob.decoder
	dec := gob.NewDecoder(f)
	err = dec.Decode(v)
	if err != nil {
		return err
	}
	return nil
}
