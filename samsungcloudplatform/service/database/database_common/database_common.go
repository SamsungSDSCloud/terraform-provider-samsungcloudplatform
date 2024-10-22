package database_common

import (
	"fmt"
	"reflect"
	"strings"
	time "time"
)

type HclKeyValueObject = map[string]interface{}
type HclSetObject = []HclKeyValueObject
type HclListObject = []interface{}

type ConvertedStruct struct {
	//BlockStorage
	BlockStorageRoleType string
	BlockStorageSize     int
	BlockStorageType     string
	BlockStorageGroupId  string

	//Server
	AvailabilityZoneName string
	ServerRoleType       string
	EpasServerName       string
	PostgresqlServerName string
	MariadbServerName    string
	SqlserverServerName  string
	MysqlServerName      string
	RedisServerName      string
	NatPublicIpId        string
	NatPublicIpAddress   string
	RedisSentinelPort    int
	CreatedDt            time.Time
}

func Contains(valueSlice []string, value string) bool {
	for _, v := range valueSlice {
		if v == value {
			return true
		}
	}
	return false
}

func ConvertObjectSliceToStructSlice(list HclListObject) []ConvertedStruct {
	result := make([]ConvertedStruct, 0)

	for _, item := range list {
		if data, ok := item.(HclKeyValueObject); ok {
			convertedStruct := ConvertedStruct{}

			for key, value := range data {
				if value == "" {
					continue
				}
				fieldName := SnakeToCamel(key)
				field := reflect.ValueOf(&convertedStruct).Elem().FieldByName(fieldName)
				switch v := value.(type) {
				case string:
					field.SetString(v)
				case int:
					field.SetInt(int64(v))
				}
			}

			result = append(result, convertedStruct)
		}
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
			case reflect.Slice:
				field.Set(reflect.ValueOf(ConvertToList(v.([]interface{}))))
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

func ConvertToList(sourceList []interface{}) []string {
	targetList := make([]string, len(sourceList))
	for i, source := range sourceList {
		targetList[i] = source.(string)
	}
	return targetList
}

func SnakeToCamel(s string) string {
	parts := strings.Split(s, "_")
	for i, p := range parts {
		parts[i] = strings.Title(strings.ToLower(p))
	}
	return strings.Join(parts, "")
}
