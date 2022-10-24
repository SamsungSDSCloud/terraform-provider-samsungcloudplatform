package common

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"regexp"
	"strconv"
	"strings"
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
	err := checkStringLength(value, 8, 20)
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
	if v.(int) < 10 {
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

func ValidateName3to28UnderscoreLowercase(v interface{}, path cty.Path) diag.Diagnostics {
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
	if !regexp.MustCompile("^[a-z][a-z\\_0-9]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
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

func ValidateName6to20(v interface{}, path cty.Path) diag.Diagnostics {
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
	if !regexp.MustCompile("[a-zA-Z0-9`~!@#^&*()-_=+|;:',.<>/?]+$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprintf("Attribute %q must contain only alpha-numerical characters", attrKey),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateCidrIpv4(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// TODO : Check CIDR for IPv4 spec and add test case

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
