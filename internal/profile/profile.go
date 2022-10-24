package profile

import (
	"errors"
	"log"
)

type Profile struct {
	Name           string
	Configurations map[string]Properties
}

func NewProfile() Profile {
	profile := Profile{}
	profile.Configurations = map[string]Properties{}
	return profile
}

func NewProfileWithName(name string) Profile {
	profile := NewProfile()
	profile.Name = name
	return profile
}

func (profile *Profile) AddProperty(name string, keyValue string) error {
	prop, ok := profile.Configurations[name]
	if !ok {
		prop = NewProperties()
		profile.Configurations[name] = prop
	}
	return prop.Add(keyValue)
}

func (profile *Profile) RemoveProperty(name string) error {
	_, ok := profile.Configurations[name]
	if !ok {
		log.Fatalln("Configuration does not exists")
		return errors.New("Configuration does not exists")
	}
	delete(profile.Configurations, name)
	return nil
}

func (profile *Profile) ClearProperty() {
	for k := range profile.Configurations {
		delete(profile.Configurations, k)
	}
}
