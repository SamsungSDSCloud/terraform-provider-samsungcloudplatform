package database_common

import (
	"github.com/SamsungSDSCloud/terraform-sdk-samsungcloudplatform/v3/library/sqlserver"
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

func TestContains(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		slice    []string
		value    string
		expected bool
	}{
		{[]string{"a", "b", "c"}, "a", true},
		{[]string{"a", "", "c"}, "o", false},
		{[]string{"aaa", "bbb", "ccc"}, "aaa", true},
	}

	for _, tc := range testCases {
		actual := Contains(tc.slice, tc.value)
		if actual != tc.expected {
			t.Error("Contains(", tc.slice, "), Values (", tc.value, ") = ", actual, ", expected ", tc.expected)
		}
	}
}

func TestConvertSecurityGroupIdList(t *testing.T) {
	t.Parallel()

	securityGroupIdList := []string{"FIREWALL-1234", "FIREWALL-5678"}
	slice := make([]interface{}, len(securityGroupIdList))
	for index, value := range securityGroupIdList {
		slice[index] = value
	}

	testCases := []struct {
		input    []interface{}
		expected []string
	}{
		{
			slice,
			[]string{"FIREWALL-1234", "FIREWALL-5678"},
		},
	}

	for _, tc := range testCases {
		actual := ConvertSecurityGroupIdList(tc.input)
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Error("ConvertSecurityGroupIdList(", tc.input, ") = (", actual, "), expected (", tc.expected, ")")
		}
	}
}

func TestMapToObjectWithCamel(t *testing.T) {
	t.Parallel()

	dnsServerIp := []string{"172.0.0.1", "172.0.0.2"}
	slice := make([]interface{}, len(dnsServerIp))
	for index, value := range dnsServerIp {
		slice[index] = value
	}

	testCases := []struct {
		input    map[string]interface{}
		expected TestModel
	}{
		{
			map[string]interface{}{
				"archive_backup_schedule_frequency": "5M",
				"backup_retention_period":           "7D",
				"backup_start_hour":                 7,
				"dns_server_ips":                    slice,
				"audit_enabled":                     true,
			},
			TestModel{"5M", "7D", 7, []string{"172.0.0.1", "172.0.0.2"}, true},
		},
		{
			map[string]interface{}{
				"archive_backup_schedule_frequency": "1M",
				"backup_retention_period":           "",
				"backup_start_hour":                 3,
				"dns_server_ips":                    slice,
				"audit_enabled":                     false,
			},
			TestModel{"1M", "", 3, []string{"172.0.0.1", "172.0.0.2"}, false},
		},
	}

	for _, tc := range testCases {
		actual := TestModel{}
		err := MapToObjectWithCamel(tc.input, &actual)
		if err != nil {
			t.Errorf("mapToObject(%v) returned an error: %v", tc.input, err)
		}
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Error("mapToObject(", tc.input, ") = (", actual, "), expected (", tc.expected, ")")
		}
	}
}

func TestConvertObjectSliceToStructSlice(t *testing.T) {
	t.Parallel()

	inputVO := map[string]interface{}{
		"block_storage_role_type": "DATA",
		"block_storage_size":      100,
		"block_storage_type":      "SSD",
		"nat_public_ip_id":        "",
	}
	objectSlice := make([]interface{}, 1)
	objectSlice[0] = inputVO

	outputVO := ConvertedStruct{
		BlockStorageRoleType: "DATA",
		BlockStorageSize:     100,
		BlockStorageType:     "SSD",
		NatPublicIpId:        "",
	}
	structSlice := make([]ConvertedStruct, 1)
	structSlice[0] = outputVO

	testCases := []struct {
		input    []interface{}
		expected []ConvertedStruct
	}{
		{
			objectSlice,
			structSlice,
		},
	}

	for _, tc := range testCases {
		actual := ConvertObjectSliceToStructSlice(tc.input)

		if !reflect.DeepEqual(actual, tc.expected) {
			t.Error("ConvertObjectSliceToStructSlice(", tc.input, ") = (", actual, "), expected (", tc.expected, ")")
		}
	}
}

func TestMapToObjectWithCamel2(t *testing.T) {
	t.Parallel()

	var dnsServerIps []interface{}
	dnsServerIps = append(dnsServerIps, "1.1.1.1")
	dnsServerIps = append(dnsServerIps, "2.2.2.2")

	expectedDnsServerIps := ConvertToList(dnsServerIps)

	testCases := []struct {
		input    map[string]interface{}
		expected sqlserver.SqlserverActiveDirectory
	}{
		{
			map[string]interface{}{
				"ad_server_user_id":       "dbuser",
				"ad_server_user_password": "1q2w",
				"dns_server_ips":          dnsServerIps,
				"domain_name":             "scp.com",
				"domain_net_bios_name":    "SCP",
				"failover_cluster_name":   "failover",
			},
			sqlserver.SqlserverActiveDirectory{"dbuser", "1q2w", expectedDnsServerIps, "scp.com", "SCP", "failover"},
		},
	}

	for _, tc := range testCases {
		actual := sqlserver.SqlserverActiveDirectory{}
		err := MapToObjectWithCamel(tc.input, &actual)
		if err != nil {
			t.Errorf("mapToObject(%v) returned an error: %v", tc.input, err)
		}
		if !reflect.DeepEqual(actual, tc.expected) {
			t.Error("mapToObject(", tc.input, ") = (", actual, "), expected (", tc.expected, ")")
		}
	}
}
