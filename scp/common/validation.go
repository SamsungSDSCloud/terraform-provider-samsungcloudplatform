package common

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

func checkStringLength(str string, min int, max int) error {
	if len(str) < min {
		return fmt.Errorf("input must be longer than %v characters", min)
	} else if len(str) > max {
		return fmt.Errorf("input must be shorter than %v characters", max)
	} else {
		return nil
	}
}

func CheckInt32Range(val int32, min int32, max int32) error {
	if val < min {
		return fmt.Errorf("the input value must be greater than or equal to %v", min)
	} else if val > max {
		return fmt.Errorf("the input value must be less than or equal to %v", max)
	} else {
		return nil
	}
}

func ValidatePassword8to20(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 8, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z\\d@$!%*#?&]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidatePassword8to30WithSpecialsExceptQuotes(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 8, 30)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z\\d@$!%*#?&(){};:\\/,.<>`~\\-\\_+=\\[\\]]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateBlockStorageSizeForOS(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	if v.(int) < 100 {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("size must be at least 10 GB"),
			AttributePath: path,
		})
	}
	if v.(int)%10 != 0 {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("size must be multiple of 10"),
			AttributePath: path,
		})
	}
	return diags
}

func ValidateBlockStorageSize(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	if v.(int) < 10 || v.(int) > 7168 {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("size must be at least 10 GB"),
			AttributePath: path,
		})
	}
	/* size must be multiple of 10 : 10 GB, 20 GB, 30GB, ...
	if v.(int)%10 != 0 {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("size must be multiple of 10"),
			AttributePath: path,
		})
	}
	*/
	return diags
}

func ValidateName3to20DashUnderscore(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z\\-\\_A-Z0-9]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName1to128DotDashUnderscore(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 1, 128)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z\\-\\_A-Z0-9\\.]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName1to256DotDashUnderscore(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 1, 256)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z\\-\\_A-Z0-9\\.]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to21LowerAlphaAndNumericWithUnderscoreStartsWithLowerAlpha(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 21)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z]([a-z0-9_]){2,20}$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only lower alpha-numeric characters and underscore", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to20NoSpecials(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z0-9]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to28Underscore(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 28)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z][a-zA-Z\\_0-9]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to20DashInMiddle(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z][a-zA-Z\\-0-9]+[a-zA-Z0-9]$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to60AlphaNumericWithSpaceDashUnderscoreStartsWithLowerAlpha(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 60)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("[a-z][a-zA-Z0-9\\-\\_\\.\\/\\s]*").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to30AlphaNumeric(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 30)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("[a-zA-Z0-9]*").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to20AlphaOnly(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName1to15AlphaOnlyStartsWithUpperCase(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 1, 15)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("[A-Z][a-zA-Z]*$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only the alphabet and begin with a capital letter", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to28AlphaDashStartsWithLowerCase(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 28)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("[a-z][a-z0-9\\-]*").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to20LowerAlphaOnly(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName2to20LowerAlphaOnly(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 2, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName2to20AlphaNumeric(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 2, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to13AlphaNumberDash(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 13)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z][a-zA-Z\\-0-9]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters with dash", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

// ^[a-zA-Z0-9+=,.@\-_ㄱ-ㅎ|ㅏ-ㅣ|가-힣]*$
func ValidateNameHangeulAlphabetSomeSpecials3to64(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	var err error = nil
	cnt := utf8.RuneCountInString(value) // cause we have hanguel here :)
	if cnt < 3 {
		err = fmt.Errorf("input must be longer than 3 characters")
	} else if cnt > 64 {
		err = fmt.Errorf("input must be shorter than 24 characters")
	}

	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z0-9+=,.@\\-_ㄱ-ㅎ|ㅏ-ㅣ|가-힣]*$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to28Dash(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 28)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z][a-zA-Z\\-0-9]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to28AlphaNumericWithSpaceAndDashStartsWithAlpha(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 28)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z][a-zA-Z\\-0-9\\ ]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName4to20NoSpecialsLowercase(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 4, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z][a-z0-9]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName6to20AlphaAndNumericWithoutSomeSpecials(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 6, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters ($ % { } [ ] " \ 제외)
	if !regexp.MustCompile("[a-zA-Z0-9`~!@#^&*()-_=+|;:',.<>/?]{6,20}$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numeric characters without some specials", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateCidrIpv4(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	_, ipNet, err := net.ParseCIDR(value)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	if ipNet.IP.To4() == nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q is an invalid IPv4 CIDR format", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateProtocol(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	const (
		Tcp   string = "TCP"
		Http  string = "HTTP"
		Https string = "HTTPS"
	)
	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	valid := strings.EqualFold(Tcp, value) || strings.EqualFold(Http, value) || strings.EqualFold(Https, value)

	if !valid {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has must be 'tcp' or 'http' or 'https'", attrKey),
			AttributePath: path,
		})
	}
	return diags
}

func ValidateSubnetType(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	const (
		PublicSubnetType  string = "PUBLIC"
		PrivateSubnetType string = "PRIVATE"
	)

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	valid := strings.EqualFold(PublicSubnetType, value) || strings.EqualFold(PrivateSubnetType, value)

	if !valid {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has must be 'public' or 'private'", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateDescriptionMaxlength50(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 0, 50)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateDescriptionMaxlength100(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 0, 100) // TODO : parameterize length
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	return diags
}

func ValidatePortRange(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := int32(v.(int))

	err := CheckInt32Range(value, 1, 65535)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateStringPortRange(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)
	port, err := strconv.Atoi(value)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Invalid port format Attribute %q has errors : %s", attrKey, value),
			AttributePath: path,
		})
		return diags
	}

	err = CheckInt32Range(int32(port), 1, 65535)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	return diags
}

func ValidatePositiveInt(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := int32(v.(int))

	err := CheckInt32Range(value, 1, 2147483647)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateIpv4(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	value := v.(string)

	trial := net.ParseIP(value)
	if trial.To4() == nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q is not IP address", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateIpv4WithEmptyValue(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	value := v.(string)

	if len(value) == 0 {
		return diags
	}

	trial := net.ParseIP(value)
	if trial.To4() == nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q is not IP address", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateThatVmStateOnlyHasRunningOrStopped(v interface{}, path cty.Path) diag.Diagnostics {

	var diags diag.Diagnostics

	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	value := v.(string)

	if !regexp.MustCompile("RUNNING|STOPPED|running|stopped").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Only RUNNING or STOPPED value of Attribute %q is allowed ", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateContractDesc(v interface{}, path cty.Path) diag.Diagnostics {

	var diags diag.Diagnostics

	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	value := v.(string)

	if !regexp.MustCompile("None|1 Year|3 Year").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Only 'None' or '1 Year' or '3 Year' value of Attribute %q is allowed ", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to20LowerAlphaAndNumberOnly(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 20)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z0-9]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only lower case letters and numbers", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateName3to28AlphaNumberDash(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get attribute key
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	// Get value
	value := v.(string)

	// Check name length
	err := checkStringLength(value, 3, 28)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q has errors : %s", attrKey, err.Error()),
			AttributePath: path,
		})
	}

	// Check characters
	if !regexp.MustCompile("^[a-z][a-z0-9-]+[a-z0-9]$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must start with a lowercase, and enter using lowercase, number, and -. However, it does not end with -.", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateServerNameInWindowImage(list HclListObject) error {

	for _, itemObject := range list {
		item := itemObject.(HclKeyValueObject)
		v, ok := item["name"]

		if !ok {
			return errors.New("There is no name attribute")
		}

		if len(v.(string)) > 15 {
			return errors.New("Servers using Windows must be 3 to 15 characters long. ")
		}

	}

	return nil
}

func ValidateContract(v interface{}, path cty.Path) diag.Diagnostics {

	var diags diag.Diagnostics

	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name

	value := v.(string)

	if !regexp.MustCompile("^(None|1 Year|3 Year)$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Only 'None' or '1 Year' or '3 Year' value of Attribute %q is allowed ", attrKey),
			AttributePath: path,
		})
	}

	return diags
}
