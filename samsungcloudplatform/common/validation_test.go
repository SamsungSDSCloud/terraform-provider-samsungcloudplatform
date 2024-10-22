package common

import (
	"github.com/hashicorp/go-cty/cty"
	"testing"
)

func TestCheckStringLength(t *testing.T) {
	if checkStringLength("", 0, 1) != nil {
		t.Error("empty should be allowed")
	}
	if checkStringLength("1", 1, 1) != nil {
		t.Error("one character should be allowed")
	}
	if checkStringLength("asdf", 5, 10) == nil {
		t.Error("shorter string should not be allowed")
	}
	if checkStringLength("asdf", 1, 3) == nil {
		t.Error("longer string shuold not be allowed")
	}
}

func TestValidateVPCName(t *testing.T) {
	ctyPath := cty.Path{
		cty.GetAttrStep{
			Name: "name",
		},
	}

	if !ValidateName3to20NoSpecials("", ctyPath).HasError() {
		t.Error("empty name should not be allowed")
	}

	if !ValidateName3to20NoSpecials("ab", ctyPath).HasError() {
		t.Error("shorter name should not be allowed")
	}

	if !ValidateName3to20NoSpecials("tooooooooooollllllllloooooooooonnnnnnnnggggggg", ctyPath).HasError() {
		t.Error("too long name should not be allowed")
	}

	if ValidateName3to20NoSpecials("asd", ctyPath).HasError() {
		t.Error("minimum length name should be allowed")
	}

	if ValidateName3to20NoSpecials("abcdefghijklmnopqrst", ctyPath).HasError() {
		t.Error("maximum length name should be allowed")
	}

	if ValidateName3to20NoSpecials("asdf123", ctyPath).HasError() {
		t.Error("alpha-numerical characters should be allowed")
	}

	if ValidateName3to20NoSpecials("123asdf", ctyPath).HasError() {
		t.Error("numerical-alpha characters should be allowed")
	}

	if !ValidateName3to20NoSpecials("asdf1!", ctyPath).HasError() {
		t.Error("special characters are not allowed")
	}

	if !ValidateName3to20NoSpecials("한글이름", ctyPath).HasError() {
		t.Error("non alpha-numercial unicode characters are not allowed")
	}

	if !ValidateName3to20NoSpecials("asd_111", ctyPath).HasError() {
		t.Error("special characters are not allowed")
	}
}

func TestValidateVPCDescription(t *testing.T) {
	ctyPath := cty.Path{
		cty.GetAttrStep{
			Name: "name",
		},
	}

	if ValidateDescriptionMaxlength50("", ctyPath).HasError() {
		t.Error("empty name should be allowed")
	}

	if !ValidateDescriptionMaxlength50("abc hijklmnopq,,.\nzyz0123456!#$789한글 띄어쓰고 098765", ctyPath).HasError() {
		t.Error("under 50 characters are allowed")
	}

	if !ValidateDescriptionMaxlength50("vvvvvvvvvvvveeeeeeeeeeeeeerrrrrrrrrrrrrrrrrryyyyyyyyyyyyyyyylllllllllllllllooooooooooooooonnnnnnnnnnnnnnnggggggggggggggggggggggggggggggg                      is  not allowed", ctyPath).HasError() {
		t.Error("over 50 characters are not allowed")
	}
}

func TestValidateName1to15AlphaOnlyStartsWithUpperCase(t *testing.T) {
	t.Parallel()

	path := cty.Path{cty.GetAttrStep{Name: "MS SQL Server DB Service Name Validation"}}

	validName := []string{
		"M",
		"MSsql",
		"MSSQL",
		"Abcdefghijklmno",
	}

	for _, v := range validName {
		if ValidateName1to15AlphaOnlyStartsWithUpperCase(v, path).HasError() {
			t.Error(path, " valid name: ", v)
		}
	}

	invalidName := []string{
		"m",
		"mssql",
		"abcdefghijklmnop",
		"abc#$%",
	}

	for _, v := range invalidName {
		if !ValidateName1to15AlphaOnlyStartsWithUpperCase(v, path).HasError() {
			t.Error(path, " invalid name: ", v)
		}
	}
}
