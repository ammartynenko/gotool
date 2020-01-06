package databox

import "errors"

const (
	errorNotFoundSection = "not found section"
	errorNotFoundKey     = "not found key"
)

type DataSection map[string]interface{}
type Databox2 map[string]DataSection

func NewDatabox2() Databox2 {
	return Databox2{}
}
func (d Databox2) Put(section, key string, value interface{}) {
	if _, found := d[section]; found {
		d[section][key] = value
	} else {
		d[section] = DataSection{}
		d[section][key] = value
	}
}
func (d Databox2) Get(section, key string) (interface{}, error) {
	if _, found := d[section]; found {
		if value, found := d[section][key]; found {
			return value, nil
		} else {
			return nil, errors.New(errorNotFoundKey)
		}
	} else {
		return nil, errors.New(errorNotFoundSection)
	}
}
func (d Databox2) GetSection(section string) (DataSection, error) {
	if value, found := d[section]; found {
		return value, nil
	}
	return nil, errors.New(errorNotFoundSection)
}
