package profile

import (
	"errors"
	"strings"
)

type Properties struct {
	Name string
	Data map[string]string
}

func NewProperties() Properties {
	properties := Properties{}
	properties.Data = map[string]string{}
	return properties
}

func NewPropertiesWithName(name string) Properties {
	properties := NewProperties()
	properties.Name = name
	return properties
}

func (properties *Properties) Add(keyValue string) error {
	slice := strings.Split(keyValue, "=")
	if len(slice) != 2 {
		return errors.New("Invalid input data")
	}
	properties.Data[slice[0]] = slice[1]
	return nil
}
func (properties *Properties) AddKeyValue(key string, value string) {
	properties.Data[key] = value
}

func (properties *Properties) Remove(key string) {
	delete(properties.Data, key)
}

func (properties *Properties) Clear() {
	for k := range properties.Data {
		delete(properties.Data, k)
	}
}
