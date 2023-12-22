package autoscaling_common

import (
	"github.com/SamsungSDSCloud/terraform-provider-samsungcloudplatform/v3/scp/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"reflect"
	"strings"
	"time"
)

func SetResponseToResourceData(response interface{}, rd *schema.ResourceData, ignoreValues ...string) {
	value := reflect.ValueOf(response)
	for i := 0; i < value.NumField(); i++ {
		if containsUsingStringSlice(ignoreValues, value.Type().Field(i).Name) {
			continue
		}
		fieldName := common.ToSnakeCase(value.Type().Field(i).Name)
		fieldValue := value.Field(i).Interface()
		if err := rd.Set(fieldName, ConvertFieldValue(fieldValue)); err != nil {
			log.Println("An error occurred in SetResponseToResourceData : ", err)
		}
	}
}

func containsUsingStringSlice(StringSlice []string, target string) bool {
	for _, s := range StringSlice {
		if strings.Contains(s, target) {
			return true
		}
	}
	return false
}

func ConvertFieldValue(fieldValue interface{}) interface{} {
	reflectValue := reflect.ValueOf(fieldValue)
	if reflectValue.Kind() == reflect.Ptr {
		// In case the type is pointer, dereference the pointer to retrieve the struct.
		for reflectValue.Kind() == reflect.Ptr {
			reflectValue = reflectValue.Elem()
		}
	}

	switch {
	case reflect.TypeOf(fieldValue) == reflect.TypeOf(time.Time{}):
		// In case the type is 'time.Time', convert to RFC3339 format.
		return fieldValue.(time.Time).Format(time.RFC3339)
	case reflectValue.Kind() == reflect.Struct:
		// In case the type is struct.
		return append(make([]interface{}, 0), ConvertStructToMap(reflectValue))
	case reflectValue.Kind() != reflect.Slice:
		// In case it's not a slice.
		return fieldValue
	case reflectValue.Len() > 0 && reflectValue.Index(0).Kind() == reflect.Struct:
		// In case of each element in the slice being a struct.
		return ConvertToStructSlice(reflectValue)
	case reflectValue.Len() > 0:
		// In case of each element in the slice being a primitive data type.
		return ConvertToPrimitiveSlice(reflectValue)
	default:
		log.Println("An unexpected case in ConvertFieldValue called by SetResponseToResourceData.")
		log.Println("type	: ", reflect.TypeOf(fieldValue))
		log.Println("value	: ", fieldValue)
		return fieldValue
	}
}

func ConvertToStructSlice(value reflect.Value) interface{} {
	var result []map[string]interface{}
	for i := 0; i < value.Len(); i++ {
		itemValue := value.Index(i)
		result = append(result, ConvertStructToMap(itemValue))
	}
	return result
}

func ConvertToPrimitiveSlice(value reflect.Value) []interface{} {
	var result []interface{}
	for i := 0; i < value.Len(); i++ {
		result = append(result, value.Index(i).Interface())
	}
	return result
}

func ConvertStructToMap(itemValue reflect.Value) map[string]interface{} {
	itemMap := make(map[string]interface{})
	for i := 0; i < itemValue.NumField(); i++ {
		fieldName := common.ToSnakeCase(itemValue.Type().Field(i).Name)
		fieldValue := itemValue.Field(i).Interface()
		if reflect.ValueOf(fieldValue).Kind() == reflect.Slice {
			fieldValue = ConvertFieldValue(fieldValue)
		}
		itemMap[fieldName] = fieldValue
	}
	return itemMap
}

func CalculateSetDifference(setA, setB *schema.Set) *schema.Set {
	difference := schema.NewSet(schema.HashString, nil)
	for _, elem := range setA.List() {
		if !setB.Contains(elem) {
			difference.Add(elem)
		}
	}
	return difference
}
