package profile

type Profiles struct {
	Name       string
	ProfileMap map[string]Profile
}

func NewProfiles() Profiles {
	profiles := Profiles{}
	profiles.ProfileMap = map[string]Profile{}
	return profiles
}

func NewProfilesWithName(name string) Profiles {
	profiles := NewProfiles()
	profiles.Name = name
	return profiles
}
