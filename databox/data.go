package databox

type Data struct {
	stock map[string]map[string]interface{}
}

func NewData() *Data {
	d := &Data{
		stock: make(map[string]map[string]interface{}),
	}
	return d
}
func (d *Data) Put(section, key string, value interface{}) {
	_, exists := d.stock[section]
	if exists {
		_, found := d.stock[section][key]
		if !found {
			d.stock[section][key] = make(map[string]interface{})
		}
	} else {
		d.stock[section] = make(map[string]interface{})
		d.stock[section][key] = make(map[string]interface{})
	}
	d.stock[section][key] = value
}
func (d *Data) Get(section, key string) (interface{}) {
	_, exists := d.stock[section]
	if exists {
		v, found := d.stock[section][key]
		if found {
			return v
		}
		return nil
	}
	return nil
}
