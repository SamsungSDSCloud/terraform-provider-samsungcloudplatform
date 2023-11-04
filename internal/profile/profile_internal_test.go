package profile_test

import (
	"testing"

	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/internal/profile"
)

func TestProfile_Empty(t *testing.T) {
	profile := profile.NewProfile()
	_, ok := profile.Configurations["a"]
	if ok {
		t.Error("Configurations must be empty")
	}
}

func TestProfile_AddProperty(t *testing.T) {
	profile := profile.NewProfile()
	profile.AddProperty("default", "hello=world")
	p, ok := profile.Configurations["default"]
	if !ok {
		t.Error("'default' property must present")
	}

	v, ok := p.Data["hello"]
	if !ok {
		t.Error("'hello' property must present")
	}
	if v != "world" {
		t.Error("'hello' property must be 'world'")
	}
}

func TestProfile_RemoveProperty(t *testing.T) {
	profile := profile.NewProfile()
	profile.AddProperty("default", "hello=world")
	profile.AddProperty("nondefault", "hello2=world2")
	_, ok := profile.Configurations["default"]
	if !ok {
		t.Error("'default' property must present")
	}
	_, ok = profile.Configurations["nondefault"]
	if !ok {
		t.Error("'nondefault' property must present")
	}

	profile.RemoveProperty("default")

	_, ok = profile.Configurations["default"]
	if ok {
		t.Error("'default' property must not present")
	}

	if len(profile.Configurations) == 0 {
		t.Error("Configurations property must not be empty")
	}

	_, ok = profile.Configurations["nondefault"]
	if !ok {
		t.Error("'nondefault' property must present")
	}
}

func TestProfile_ClearProperty(t *testing.T) {
	profile := profile.NewProfile()
	profile.AddProperty("default", "hello=world")
	profile.AddProperty("nondefault", "hello2=world2")
	_, ok := profile.Configurations["default"]
	if !ok {
		t.Error("'default' property must present")
	}
	_, ok = profile.Configurations["nondefault"]
	if !ok {
		t.Error("'default' property must present")
	}

	profile.ClearProperty()

	if len(profile.Configurations) != 0 {
		t.Error("Configurations property must be empty")
	}
}
