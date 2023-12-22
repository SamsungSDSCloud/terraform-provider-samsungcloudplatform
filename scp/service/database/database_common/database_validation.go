package database_common

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
)

func getAttrKey(path cty.Path) string {
	attr := path[len(path)-1].(cty.GetAttrStep)
	attrKey := attr.Name
	return attrKey
}

func ValidServerState(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attrKey := getAttrKey(path)
	value := v.(string)

	if !regexp.MustCompile("^(RUNNING|STOPPED)$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprint("Only RUNNING or STOPPED value of Attribute ", attrKey, " is allowed : ", value),
			AttributePath: path,
		})
	}
	return diags
}

func ValidateContractPeriod(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attrKey := getAttrKey(path)
	value := v.(string)

	if !regexp.MustCompile("^(None|1 Year|3 Year)$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprint("Only 'None', '1 Year' or '3 Year' value of Attribute ", attrKey, " is allowed : ", value),
			AttributePath: path,
		})
	}
	return diags
}

func ValidateBackupScheduleFrequency(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attrKey := getAttrKey(path)
	value := v.(string)

	if !regexp.MustCompile("^(5M|10M|30M|1H)$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprint("Only '5M', '10M', '30M', '1H' value of Attribute ", attrKey, " is allowed : ", value),
			AttributePath: path,
		})
	}
	return diags
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

func ValidateBackupStartHour(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attrKey := getAttrKey(path)
	value := v.(int)

	if value < 0 || value > 23 {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprint("Only between 0 to 23 value of Attribute ", attrKey, " is allowed : ", value),
			AttributePath: path,
		})
	}

	return diags
}

func ValidateSameValue(expectedValue string) schema.SchemaValidateDiagFunc {
	return func(v interface{}, path cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics

		attrKey := getAttrKey(path)
		value := v.(string)

		if v.(string) != expectedValue {
			diags = append(diags, diag.Diagnostic{
				Severity:      diag.Error,
				Summary:       fmt.Sprint("In case of ", attrKey, ", only \"", expectedValue, "\" is allowed : ", value),
				AttributePath: path,
			})
		}

		return diags
	}
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

func ValidateBlockStorageType(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attrKey := getAttrKey(path)
	value := v.(string)

	if !regexp.MustCompile("^(SSD|HDD)$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprint("Only SSD or HDD value of Attribute ", attrKey, " is allowed : ", value),
			AttributePath: path,
		})
	}
	return diags
}

func ValidateBlockStorageRoleType(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attrKey := getAttrKey(path)
	value := v.(string)

	if !regexp.MustCompile("^(DATA|ARCHIVE|TEMP|BACKUP)$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprint("Only DATA, ARCHIVE, TEMP and BACKUP value of Attribute ", attrKey, " is allowed : ", value),
			AttributePath: path,
		})
	}
	return diags
}

func ValidateBlockStorageSize(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics
	value := v.(int)

	if value < 10 || value > 5120 {
		return diag.Errorf("Port Number must be between 1024 and 65535")
	}

	return diags
}

func ValidateAvailabilityZone(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attrKey := getAttrKey(path)
	value := v.(string)

	if !regexp.MustCompile("^(AZ1|AZ2)$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprint("Only AZ1 or AZ2 value of Attribute ", attrKey, " is allowed : ", value),
			AttributePath: path,
		})
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

func ValidateServerRoleType(v interface{}, path cty.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	attrKey := getAttrKey(path)
	value := v.(string)

	if !regexp.MustCompile("^(ACTIVE|STANDBY)$").MatchString(value) {
		diags = append(diags, diag.Diagnostic{
			Severity:      diag.Error,
			Summary:       fmt.Sprint("Only ACTIVE or STANDBY value of Attribute ", attrKey, " is allowed : ", value),
			AttributePath: path,
		})
	}
	return diags
}
