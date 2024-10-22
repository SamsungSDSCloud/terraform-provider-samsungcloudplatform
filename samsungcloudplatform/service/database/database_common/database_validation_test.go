package database_common

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initPath(v string) cty.Path {
	path := cty.Path{cty.GetAttrStep{Name: v}}
	return path
}

func TestValidateIntegerInRange(t *testing.T) {
	validateFunc := ValidateIntegerInRange(10, 20)

	// Define test cases
	tests := []struct {
		name          string
		input         int
		expectedError bool
	}{
		{"valid value", 10, false},
		{"valid value", 20, false},
		{"invalid value", 9, true},
		{"invalid value", 21, true},
		{"invalid value", 0, true},
		{"invalid value", -10, true},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := validateFunc(tc.input, initPath("testValidateIntegerInRange"))

			if tc.expectedError {
				assert.NotEmpty(t, diags, "Expected an error but got none")
			} else {
				assert.Empty(t, diags, "Expected no error but got some")
			}
		})
	}
}

func TestValidateIntegerGreaterEqualThan(t *testing.T) {
	validateFunc := ValidateIntegerGreaterEqualThan(10)

	// Define test cases
	tests := []struct {
		name          string
		input         int
		expectedError bool
	}{
		{"valid value", 10, false},
		{"valid value", 15, false},
		{"invalid value", 9, true},
		{"invalid value", 0, true},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := validateFunc(tc.input, initPath("testValidateIntegerGreaterEqualThan"))

			if tc.expectedError {
				assert.NotEmpty(t, diags, "Expected an error but got none")
			} else {
				assert.Empty(t, diags, "Expected no error but got some")
			}
		})
	}
}

func TestValidateIntegerLessEqualThan(t *testing.T) {
	validateFunc := ValidateIntegerLessEqualThan(10)

	// Define test cases
	tests := []struct {
		name          string
		input         int
		expectedError bool
	}{
		{"valid value", 10, false},
		{"valid value", 5, false},
		{"invalid value", 11, true},
		{"invalid value", 10000, true},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := validateFunc(tc.input, initPath("testValidateIntegerGreaterEqualThan"))

			if tc.expectedError {
				assert.NotEmpty(t, diags, "Expected an error but got none")
			} else {
				assert.Empty(t, diags, "Expected no error but got some")
			}
		})
	}
}

func TestValidateStringInOptions(t *testing.T) {
	validateFunc := ValidateStringInOptions("option1", "option2", "option3")

	// Define test cases
	tests := []struct {
		name          string
		input         string
		expectedError bool
	}{
		{"valid value", "option1", false},
		{"valid value", "option2", false},
		{"valid value", "option3", false},
		{"invalid value", "option", true},
		{"invalid value", "", true},
		{"invalid value", "option33", true},
		{"invalid value", "AZ1", true},
	}

	// Run tests
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			diags := validateFunc(tc.input, initPath("testStringInOption"))

			if tc.expectedError {
				assert.NotEmpty(t, diags, "Expected an error but got none")
			} else {
				assert.Empty(t, diags, "Expected no error but got some")
			}
		})
	}
}

func TestValidateBackupRetentionPeriod(t *testing.T) {
	t.Parallel()

	ctyPath := initPath("Backup retention Period")

	validState := []string{
		"7D",
		"14D",
		"25D",
		"35D",
	}

	for _, v := range validState {
		if ValidateBackupRetentionPeriod(v, ctyPath).HasError() {
			t.Error(ctyPath, " should be a valid state: ", v)
		}
	}

	invalidState := []string{
		"7",
		"sample",
		"25A",
		"40D",
		"한글",
	}

	for _, v := range invalidState {
		if !ValidateBackupRetentionPeriod(v, ctyPath).HasError() {
			t.Error(ctyPath, " should be a invalid state: ", v)
		}
	}
}

func TestValidateAlphaNumeric3to20(t *testing.T) {
	t.Parallel()

	ctyPath := initPath("database name")

	validState := []string{
		"abc",
		"32lkaelkwe",
		"09090",
		"ABCDEFGHIJ",
	}

	for _, v := range validState {
		if ValidateAlphaNumeric3to20(v, ctyPath).HasError() {
			t.Error(ctyPath, " should be a valid state: ", v)
		}
	}

	invalidState := []string{
		"Test_sample",
		"abcdefghijklmnopqrstuvwxyz",
		"ab",
		"0999#",
	}

	for _, v := range invalidState {
		if !ValidateAlphaNumeric3to20(v, ctyPath).HasError() {
			t.Error(ctyPath, " should be a invalid state: ", v)
		}
	}
}

func TestValidatePortNumber(t *testing.T) {
	t.Parallel()

	ctyPath := initPath("Port number")

	validState := []int{
		1024,
		65535,
	}

	for _, v := range validState {
		if ValidatePortNumber(v, ctyPath).HasError() {
			t.Error(ctyPath, " should be a valid state: ", v)
		}
	}

	invalidState := []int{
		1023,
		65536,
	}

	for _, v := range invalidState {
		if !ValidatePortNumber(v, ctyPath).HasError() {
			t.Error(ctyPath, " should be a invalid state: ", v)
		}
	}
}

func TestValidateBlockStorageSize(t *testing.T) {
	t.Parallel()

	ctyPath := initPath("Block Storage Size")

	validState := []int{
		10,
		5120,
	}

	for _, v := range validState {
		if ValidateBlockStorageSize(v, ctyPath).HasError() {
			t.Error(ctyPath, " should be a valid state: ", v)
		}
	}

	invalidState := []int{
		9,
		5121,
	}

	for _, v := range invalidState {
		if !ValidateBlockStorageSize(v, ctyPath).HasError() {
			t.Error(ctyPath, " should be a invalid state: ", v)
		}
	}
}

func TestValidate3to20LowercaseNumberDashAndStartLowercase(t *testing.T) {
	t.Parallel()

	path := initPath("Block Storage Type")

	validState := []string{
		"sampleserver1",
		"sampleservsampleserv",
		"abc",
		"a00000-------",
	}

	for _, v := range validState {
		if Validate3to20LowercaseNumberDashAndStartLowercase(v, path).HasError() {
			t.Error(path, " should be a valid state: ", v)
		}
	}

	invalidState := []string{
		"ab",
		"0sampleserver",
		"Sampleserver",
		"sampleservsampleserv1",
		"abc#$%",
	}

	for _, v := range invalidState {
		if !Validate3to20LowercaseNumberDashAndStartLowercase(v, path).HasError() {
			t.Error(path, " should be a invalid state: ", v)
		}
	}
}

func TestValidate3to15LowercaseNumberDashAndStartLowercase(t *testing.T) {
	t.Parallel()

	path := initPath("Block Storage Type")

	validState := []string{
		"sampleserver1",
		"sampleservsampl",
		"abc",
		"a00000-------",
	}

	for _, v := range validState {
		if Validate3to15LowercaseNumberDashAndStartLowercase(v, path).HasError() {
			t.Error(path, " should be a valid state: ", v)
		}
	}

	invalidState := []string{
		"ab",
		"0sampleserver",
		"Sampleserver",
		"sampleservsample",
		"abc#$%",
	}

	for _, v := range invalidState {
		if !Validate3to15LowercaseNumberDashAndStartLowercase(v, path).HasError() {
			t.Error(path, " should be a invalid state: ", v)
		}
	}
}

func TestValidateDbUserName2to20AlphaNumeric(t *testing.T) {
	t.Parallel()

	path := initPath("MS SQL Server Database Username")

	validState := []string{
		"123",
		"db",
		"dbuser",
		"dbuser123",
		"db123user",
		"db_user",
		"db-user",
		"dbuser-_",
		"DB_user",
	}

	for _, v := range validState {
		if ValidateDbUserName2to20AlphaNumeric(v, path).HasError() {
			t.Error(path, " should be a valid Username : ", v)
		}
	}

	invalidState := []string{
		"d",
		"db_user12345678901234567890",
		"db-user!",
		"@dbuser-_",
	}

	for _, v := range invalidState {
		if !ValidateDbUserName2to20AlphaNumeric(v, path).HasError() {
			t.Error(path, " should be a invalid Username: ", v)
		}
	}
}

func TestValidatePassword8to30WithSpecialsExceptQuotesAndDollar(t *testing.T) {
	t.Parallel()

	path := initPath("MS SQL Server Database Password")

	validState := []string{
		"abc123456!!",
		"dlqtl#00",
		"abcd!@#1987",
	}

	for _, v := range validState {
		if ValidatePassword8to30WithSpecialsExceptQuotesAndDollar(v, path).HasError() {
			t.Error(path, " should be a valid Password : ", v)
		}
	}

	invalidState := []string{
		"abc",
		"1234",
		"abc1234",
		"1q2w3e4r5t",
		"1234567890",
		"123456789012345678901234567890abc",
		"abcdefghij",
		"abcd1234",
		"dlatl000",
		"dlqtl$00",
		"dlqtl#$00",
		"dlqtl$#00",
		"dlqtl'#00",
		"dlqtl\"#00",
		"'dlatl#00",
		"\"dlatl#00",
		"dlatl#00'",
		"dlatl#00\"",
	}

	for _, v := range invalidState {
		if !ValidatePassword8to30WithSpecialsExceptQuotesAndDollar(v, path).HasError() {
			t.Error(path, " should be a invalid Password: ", v)
		}
	}
}
