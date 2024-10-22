package database_common

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
	"strings"
)

func getAttrKey(path cty.Path) string {
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name
	return attrKey
}

func ValidateIntegerInRange(min, max int) schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		attrKey := getAttrKey(path)
		value := v.(int)

		if value < min || value > max {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprint("In case of ", attrKey, ", Value must be between ", min, " and ", max, ", but god ", value),
				AttributePath: path,
			})
		}

		return diags
	}
}

func ValidateIntegerGreaterEqualThan(min int) schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		attrKey := getAttrKey(path)
		value := v.(int)

		if value < min {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprint("In case of ", attrKey, ", value must be greater than or equal to  ", min, ", but god ", value),
				AttributePath: path,
			})
		}

		return diags
	}
}

func ValidateIntegerLessEqualThan(max int) schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		attrKey := getAttrKey(path)
		value := v.(int)

		if value > max {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprint("In case of ", attrKey, ", value must be less than or equal to  ", max, ", but god ", value),
				AttributePath: path,
			})
		}

		return diags
	}
}

func ValidateStringInOptions(options ...string) schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		attrKey := getAttrKey(path)
		value := v.(string)
		found := false

		for _, option := range options {
			if v.(string) == option {
				found = true
				break
			}
		}

		if !found {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprint("In case of ", attrKey, ", \"", strings.Join(options, "\", \""), "\" is allowed : ", value),
				AttributePath: path,
			})
		}

		return diags
	}
}

func ValidateBackupRetentionPeriod(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attrKey := getAttrKey(path)
	value := v.(string)

	if !regexp.MustCompile("^(7|8|9|[1-2][0-9]|3[0-5])D$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprint("Only between 7d to 35d value of Attribute ", attrKey, " is allowed : ", value),
			AttributePath: path,
		})
	}
	return diags
}
func ValidateAlphaNumeric3to20(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	value := v.(string)

	if len(value) < 3 || len(value) > 20 {
		return diag.Errorf("Field value must be between 3 and 20 characters")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(value) {
		return diag.Errorf("Field value must contain only letters and numbers")
	}

	return diags
}

func ValidatePortNumber(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	value := v.(int)

	if value < 1024 || value > 65535 {
		return diag.Errorf("Port Number must be between 1024 and 65535")
	}

	return diags
}

func ValidateBlockStorageSize(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	value := v.(int)

	if value < 10 || value > 5120 {
		return diag.Errorf("BlockStorage Size must be between 10 and 5120")
	}

	return diags
}

func Validate3to20LowercaseNumberDashAndStartLowercase(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	value := v.(string)

	if len(value) < 3 || len(value) > 20 {
		return diag.Errorf("Field value must be between 3 and 20 characters")
	}

	if !regexp.MustCompile(`^[a-z][a-z0-9\\-]+$`).MatchString(value) {
		return diag.Errorf("It should consist of lowercase letters, numbers, and dashes, and the first character should be an lowercase letter.")
	}

	return diags
}

func Validate3to15LowercaseNumberDashAndStartLowercase(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	value := v.(string)

	if len(value) < 3 || len(value) > 15 {
		return diag.Errorf("Field value must be between 3 and 15 characters")
	}

	if !regexp.MustCompile(`^[a-z][a-z0-9\\-]+$`).MatchString(value) {
		return diag.Errorf("It should consist of lowercase letters, numbers, and dashes, and the first character should be an lowercase letter.")
	}

	return diags
}

func ValidateDbUserName2to20AlphaNumeric(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	// Get value
	value := v.(string)

	// Check name length
	if len(value) < 2 || len(value) > 20 {
		return diag.Errorf("Field value must be between 2 and 20 characters")
	}

	// Check characters
	if !regexp.MustCompile("^[a-zA-Z0-9-_]+$").MatchString(value) {
		return diag.Errorf("Field value must contain only letters and numbers")
	}

	return diags
}

func checkStringLength(str string, min int, max int) error {
	if len(str) < min {
		return fmt.Errorf("input must be longer than %v characters", min)
	} else if len(str) > max {
		return fmt.Errorf("input must be shorter than %v characters", max)
	} else {
		return nil
	}
}
func ValidatePassword8to30WithSpecialsExceptQuotesAndDollar(v interface{}, path cty.Path) diag.Diagnostics {
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

	// number
	if !regexp.MustCompile("\\d").MatchString(value) {
		return diag.Errorf("Field value must contain numbers")
	}

	// character
	if !regexp.MustCompile("^.*[a-zA-Z].*$").MatchString(value) {
		return diag.Errorf("Field value must contain letters")
	}

	// special
	//if !regexp.MustCompile("^[@!%*#?&(){};:/,.<>`~\\-_+=\\[\\]]*").MatchString(value) {
	if !regexp.MustCompile("[@!%*#?&(){};:/,.<>`~\\-_+=\\[\\]]").MatchString(value) {
		return diag.Errorf("Field value must contain special letter")
	}

	// dollar letter
	if regexp.MustCompile(".*\\$.*").MatchString(value) {
		return diag.Errorf("Field value must not contain special letter($)")
	}

	// single quotation mark
	if regexp.MustCompile(".*'.*").MatchString(value) {
		return diag.Errorf("Field value must not contain special letter(')")
	}

	// double quotation mark
	if regexp.MustCompile(".*\".*").MatchString(value) {
		return diag.Errorf("Field value must not contain special letter(\")")
	}

	return diags
}
