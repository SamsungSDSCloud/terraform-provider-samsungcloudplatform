package database_common

import (
	"github.com/hashicorp/go-cty/cty"
	"testing"
)

func initPath(v string) cty.Path {
	path := cty.Path{cty.GetAttrStep{Name: v}}
	return path
}

func TestValidServerState(t *testing.T) {
	t.Parallel()

	path := initPath("Server state")

	validState := []string{
		"RUNNING",
		"STOPPED",
	}

	for _, v := range validState {
		if ValidServerState(v, path).HasError() {
			t.Error(path, " should be a valid state: ", v)
		}
	}

	invalidState := []string{
		"abc",
		"RUNNINGS",
		"running",
		"12355",
	}

	for _, v := range invalidState {
		if !ValidServerState(v, path).HasError() {
			t.Error(path, " should be a invalid state: ", v)
		}
	}
}

func TestValidateContractPeriod(t *testing.T) {
	t.Parallel()

	ctyPath := initPath("Contract Period")

	validState := []string{
		"None",
		"1 Year",
		"3 Year",
	}

	for _, v := range validState {
		if ValidateContractPeriod(v, ctyPath).HasError() {
			t.Error(ctyPath, " should be a valid state: ", v)
		}
	}

	invalidState := []string{
		"NONE",
		"1-Year",
		"5 Year",
		"한글",
	}

	for _, v := range invalidState {
		if !ValidateContractPeriod(v, ctyPath).HasError() {
			t.Error(ctyPath, " should be a invalid state: ", v)
		}
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

func TestValidateBlockStorageType(t *testing.T) {
	t.Parallel()

	path := initPath("Block Storage Type")

	validState := []string{
		"SSD",
		"HDD",
	}

	for _, v := range validState {
		if ValidateBlockStorageType(v, path).HasError() {
			t.Error(path, " should be a valid state: ", v)
		}
	}

	invalidState := []string{
		"SDS",
		"SSSD",
		"running",
		"12355",
	}

	for _, v := range invalidState {
		if !ValidateBlockStorageType(v, path).HasError() {
			t.Error(path, " should be a invalid state: ", v)
		}
	}
}

func TestValidateBlockStorageRoleType(t *testing.T) {
	t.Parallel()

	path := initPath("Block Storage Role Type")

	validState := []string{
		"DATA",
		"ARCHIVE",
		"TEMP",
		"BACKUP",
	}

	for _, v := range validState {
		if ValidateBlockStorageRoleType(v, path).HasError() {
			t.Error(path, " should be a valid state: ", v)
		}
	}

	invalidState := []string{
		"DATE",
		"TMP",
		"Database",
		"12355",
	}

	for _, v := range invalidState {
		if !ValidateBlockStorageRoleType(v, path).HasError() {
			t.Error(path, " should be a invalid state: ", v)
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

func TestValidateAvailabilityZone(t *testing.T) {
	t.Parallel()

	path := initPath("Block Storage Type")

	validState := []string{
		"AZ1",
		"AZ2",
	}

	for _, v := range validState {
		if ValidateAvailabilityZone(v, path).HasError() {
			t.Error(path, " should be a valid state: ", v)
		}
	}

	invalidState := []string{
		"AZ3",
		"AZ0",
		"AZ",
		"12355",
	}

	for _, v := range invalidState {
		if !ValidateAvailabilityZone(v, path).HasError() {
			t.Error(path, " should be a invalid state: ", v)
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

func TestValidateServerRoleType(t *testing.T) {
	t.Parallel()

	path := initPath("Block Storage Type")

	validState := []string{
		"ACTIVE",
		"STANDBY",
	}

	for _, v := range validState {
		if ValidateServerRoleType(v, path).HasError() {
			t.Error(path, " should be a valid state: ", v)
		}
	}

	invalidState := []string{
		"TIVE",
		"ACTIVE0",
		"########",
	}

	for _, v := range invalidState {
		if !ValidateServerRoleType(v, path).HasError() {
			t.Error(path, " should be a invalid state: ", v)
		}
	}
}
