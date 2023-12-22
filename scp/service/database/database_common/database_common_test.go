package database_common

import (
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/postgresql"
	"reflect"
	"testing"
)

func TestSnakeToCamel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    string
		expected string
	}{
		{"hello_world", "HelloWorld"},
		{"snake_case_example", "SnakeCaseExample"},
		{"", ""},
		{"a", "A"},
		{"_", ""},
	}

	for _, tc := range testCases {
		actual := SnakeToCamel(tc.input)
		if actual != tc.expected {
			t.Error("SnakeToCamel(", tc.input, ") = ", actual, ", expected ", tc.expected)
		}
	}
}

func TestMapToObjectWithCamel(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input    map[string]interface{}
		expected postgresql.DbClusterCreateFullBackupConfigRequest
	}{
		{
			map[string]interface{}{
				"object_storage_id":                 "S3OBJECTSTORAGE-XXXXX",
				"archive_backup_schedule_frequency": "5M",
				"backup_retention_period":           "7D",
				"backup_start_hour":                 7,
			},
			postgresql.DbClusterCreateFullBackupConfigRequest{"5M", "7D", 7, "S3OBJECTSTORAGE-XXXXX"},
		},
		{
			map[string]interface{}{
				"object_storage_id":                 "S3OBJECTSTORAGE-YYYYY",
				"archive_backup_schedule_frequency": "1H",
				"backup_retention_period":           "15D",
				"backup_start_hour":                 3,
			},
			postgresql.DbClusterCreateFullBackupConfigRequest{"1H", "15D", 3, "S3OBJECTSTORAGE-YYYYY"},
		},
	}

	for _, tc := range testCases {
		actual := postgresql.DbClusterCreateFullBackupConfigRequest{}
		err := MapToObjectWithCamel(tc.input, &actual)
		if err != nil {
			t.Errorf("mapToObject(%v) returned an error: %v", tc.input, err)
		}
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Error("mapToObject(", tc.input, ") = (", actual, "), expected (", tc.expected, ")")
		}
	}
}
