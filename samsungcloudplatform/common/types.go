package common

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"reflect"
	"strings"
	"time"
)

// TODO : contents to resource data (with filter)

func SetResponseToResourceData(response interface{}, rd *schema.ResourceData, ignoreValues ...string) {
	value := reflect.ValueOf(response)
	for i := 0; i < value.NumField(); i++ {
		if containsUsingStringSlice(ignoreValues, value.Type().Field(i).Name) {
			continue
		}
		fieldName := ToSnakeCase(value.Type().Field(i).Name)
		fieldValue := value.Field(i).Interface()
		if err := rd.Set(fieldName, convertFieldValue(fieldValue)); err != nil {
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

func convertFieldValue(fieldValue interface{}) interface{} {
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
		return append(make([]interface{}, 0), convertStructToMap(reflectValue))
	case reflectValue.Kind() != reflect.Slice:
		// In case it's not a slice.
		return fieldValue
	case reflectValue.Len() > 0 && reflectValue.Index(0).Kind() == reflect.Struct:
		// In case of each element in the slice being a struct.
		return convertToStructSlice(reflectValue)
	case reflectValue.Len() > 0:
		// In case of each element in the slice being a primitive data type.
		return convertToPrimitiveSlice(reflectValue)
	default:
		log.Println("An unexpected case in convertFieldValue called by SetResponseToResourceData.")
		log.Println("type	: ", reflect.TypeOf(fieldValue))
		log.Println("value	: ", fieldValue)
		return fieldValue
	}
}

func convertToStructSlice(value reflect.Value) interface{} {
	var result []map[string]interface{}
	for i := 0; i < value.Len(); i++ {
		itemValue := value.Index(i)
		result = append(result, convertStructToMap(itemValue))
	}
	return result
}

func convertToPrimitiveSlice(value reflect.Value) []interface{} {
	var result []interface{}
	for i := 0; i < value.Len(); i++ {
		result = append(result, value.Index(i).Interface())
	}
	return result
}

func convertStructToMap(itemValue reflect.Value) map[string]interface{} {
	itemMap := make(map[string]interface{})
	for i := 0; i < itemValue.NumField(); i++ {
		fieldName := ToSnakeCase(itemValue.Type().Field(i).Name)
		fieldValue := itemValue.Field(i).Interface()
		if reflect.ValueOf(fieldValue).Kind() == reflect.Slice {
			fieldValue = convertFieldValue(fieldValue)
		}
		itemMap[fieldName] = fieldValue
	}
	return itemMap
}
