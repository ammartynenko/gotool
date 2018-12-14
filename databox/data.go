package databox

type Data struct {
	Stocker map[string]map[string]interface{}
}

func NewData() *Data {
	d := &Data{
		Stocker: make(map[string]map[string]interface{}),
	}
	return d
}
func (d *Data) Put(section, key string, value interface{}) {
	_, exists := d.Stocker[section]
	if exists {
		_, found := d.Stocker[section][key]
		if !found {
			d.Stocker[section][key] = make(map[string]interface{})
		}
	} else {
		d.Stocker[section] = make(map[string]interface{})
		d.Stocker[section][key] = make(map[string]interface{})
	}
	d.Stocker[section][key] = value
}
func (d *Data) Get(section, key string) (interface{}) {
	_, exists := d.Stocker[section]
	if exists {
		v, found := d.Stocker[section][key]
		if found {
			return v
		}
		return nil
	}
	return nil
}
