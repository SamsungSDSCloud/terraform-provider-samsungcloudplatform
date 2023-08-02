package profile_test

import (
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v2/internal/profile"
)

func TestProperties_Empty(t *testing.T) {
	properties := profile.NewProperties()
	_, ok := properties.Data["a"]
	if ok {
		t.Error("Data must be empty")
	}
}

func TestProperties_Add(t *testing.T) {
	properties := profile.NewProperties()
	properties.Add("hello=world!")

	v, ok := properties.Data["hello"]
	if !ok {
		t.Error("Key not found")
	}
	if v != "world!" {
		t.Error("Value is invalid")
	}
}

func TestProperties_AddKeyValue(t *testing.T) {
	properties := profile.NewProperties()
	properties.AddKeyValue("hello", "world!")

	v, ok := properties.Data["hello"]
	if !ok {
		t.Error("Key not found")
	}
	if v != "world!" {
		t.Error("Value is invalid")
	}
}

func TestProperties_Remove(t *testing.T) {
	properties := profile.NewProperties()
	properties.Add("hello=world!")
	properties.Add("hello2=world2!")
	_, ok := properties.Data["hello"]
	if !ok {
		t.Error("Key 'hello' must exist")
	}
	properties.Remove("hello")
	_, ok = properties.Data["hello"]
	if ok {
		t.Error("Key 'hello' must not exist")
	}
	if len(properties.Data) == 0 {
		t.Error("Data must not be empty")
	}
	_, ok = properties.Data["hello2"]
	if !ok {
		t.Error("Key 'hello2' must exist")
	}
}

func TestProperties_Clear(t *testing.T) {
	properties := profile.NewProperties()
	properties.Add("hello=world!")
	properties.Add("hello2=world!")
	_, ok := properties.Data["hello"]
	if !ok {
		t.Error("Key 'hello' must exist")
	}
	_, ok = properties.Data["hello2"]
	if !ok {
		t.Error("Key 'hello2' must exist")
	}
	properties.Clear()
	if len(properties.Data) != 0 {
		t.Error("Data must be empty")
	}
}
