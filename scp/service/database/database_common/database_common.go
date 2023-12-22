package database_common

import (
	"fmt"
	"reflect"
	"strings"
)

type HclKeyValueObject = map[string]interface{}
type HclSetObject = []HclKeyValueObject
type HclListObject = []interface{}

type BlockStorage struct {
	BlockStorageRoleType string
	BlockStorageSize     int
	BlockStorageType     string
}

type PostgresqlServer struct {
	AvailabilityZoneName string
	PostgresqlServerName string
	ServerRoleType       string
}

func Contains(valueSlice []string, value string) bool {
	for _, v := range valueSlice {
		if v == value {
			return true
		}
	}
	return false
}

//TODO convert 함수를 하나로 합칠수 있을지 고민
func ConvertBlockStorageList(list HclListObject) []BlockStorage {
	var result []BlockStorage
	for _, itemObject := range list {
		item := itemObject.(HclKeyValueObject)
		var info BlockStorage
		if v, ok := item["block_storage_role_type"]; ok {
			info.BlockStorageRoleType = v.(string)
		}
		if v, ok := item["block_storage_size"]; ok {
			info.BlockStorageSize = v.(int)
		}
		if v, ok := item["block_storage_type"]; ok {
			info.BlockStorageType = v.(string)
		}
		result = append(result, info)
	}
	return result
}

func ConvertServerList(list HclListObject) []PostgresqlServer {
	var result []PostgresqlServer
	for _, itemObject := range list {
		item := itemObject.(HclKeyValueObject)
		var info PostgresqlServer
		if v, ok := item["availability_zone_name"]; ok {
			info.AvailabilityZoneName = v.(string)
		}
		if v, ok := item["postgresql_server_name"]; ok {
			info.PostgresqlServerName = v.(string)
		}
		if v, ok := item["server_role_type"]; ok {
			info.ServerRoleType = v.(string)
		}
		result = append(result, info)
	}
	return result
}

func MapToObjectWithCamel(m map[string]interface{}, obj interface{}) error {
	for k, v := range m {
		key := SnakeToCamel(k)
		field := reflect.ValueOf(obj).Elem().FieldByName(key)
		if field.IsValid() {
			fmt.Println(v)
			switch field.Kind() {
			case reflect.String:
				field.SetString(v.(string))
			case reflect.Int32:
				field.SetInt(int64(v.(int)))
			case reflect.Bool:
				field.SetBool(v.(bool))
			}
		}
	}

	return nil
}

func ConvertSecurityGroupIdList(securityGroupIds []interface{}) []string {
	securityGroupIdList := make([]string, len(securityGroupIds))
	for i, valueIpv4 := range securityGroupIds {
		securityGroupIdList[i] = valueIpv4.(string)
	}
	return securityGroupIdList
}

func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		parts[i] = strings.Title(strings.ToLower(p))
	}
	return strings.Join(parts, "")
}
