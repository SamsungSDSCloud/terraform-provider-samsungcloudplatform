package common

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func DatasourceFilter() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Required:    true,
					Description: "Filtering target name",
				},
				"values": {
					Type:        schema.TypeList,
					Required:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Filtering values. Each matching value is appended. (OR rule)",
				},
				"use_regex": {
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
					Description: "Enable regex match for values",
				},
			},
		},
	}
}

func ApplyFilter(schemaMap map[string]*schema.Schema, filter *schema.Set, items HclSetObject) HclSetObject {
	if filter == nil || filter.Len() == 0 {
		return items
	}

	for _, f := range filter.List() {
		kv := f.(HclKeyValueObject)
		targetName := kv["name"].(string)

		var err error
		var elements []string
		if elements, err = getElements(schemaMap, targetName); err != nil {
			// Fallback to default
			elements = []string{targetName}
		}

		useRegex := false
		if r, regexOk := kv["use_regex"]; regexOk {
			useRegex = r.(bool)
		}

		// create a string equality check strategy based on this filters "regex" flag
		stringsEqual := func(propertyVal string, filterVal string) bool {
			if useRegex {
				re, err := regexp.Compile(filterVal)
				if err != nil {
					log.Printf(`[WARN] Invalid regular expression "%s" for "%s" filter\n`, filterVal, targetName)
					return false
				}
				return re.MatchString(propertyVal)
			}

			return filterVal == propertyVal
		}

		result := make(HclSetObject, 0)
		for _, item := range items {
			targetVal, targetValOk := getValueFromPath(item, elements)
			if targetValOk && orComparator(targetVal, kv["values"].([]interface{}), stringsEqual) {
				result = append(result, item)
			}
		}
		items = result
	}
	return items
}

func getElements(schemaMap map[string]*schema.Schema, targetName string) ([]string, error) {
	if schemaMap == nil {
		return nil, fmt.Errorf("input schema is nil")
	}

	tokenized := strings.Split(targetName, ".")

	if len(tokenized) == 0 {
		return nil, fmt.Errorf("invalid target name : %s", targetName)
	}

	var elements []string
	currentSchema := schemaMap
	for i, t := range tokenized {
		if fieldSchema, ok := currentSchema[t]; ok && checkValidSchema(fieldSchema) {
			// Add first
			elements = append(elements, t)

			// Check for nested items
			convertedElementSchema, conversionOk := fieldSchema.Elem.(*schema.Resource)
			if !conversionOk {
				if len(tokenized) > i+1 {
					if fieldSchema.Type != schema.TypeMap {
						return nil, fmt.Errorf("invalid target name format found %s", targetName)

					}
					e := strings.Join(tokenized[i+1:], ".")
					elements = append(elements, e)
				}
				break
			} else {
				// get next schema and handle next token
				currentSchema = convertedElementSchema.Schema
			}
		} else {
			return nil, fmt.Errorf("invalid schema found for filter name %s", targetName)
		}
	}

	if len(elements) == 0 {
		return nil, fmt.Errorf("everything is filtered out")
	}

	return elements, nil
}

func checkValidSchema(s *schema.Schema) bool {
	if s.Type == schema.TypeList || s.Type == schema.TypeSet {
		if elemSchema, conversionOk := s.Elem.(*schema.Schema); conversionOk && elemSchema.Type == schema.TypeString {
			return true
		} else if s.MaxItems == 1 && s.MinItems == 1 {
			return true
		}
		return false
	}
	return true
}

type StringCheck func(propertyVal string, filterVal string) bool

func orComparator(target interface{}, filters []interface{}, stringsEqual StringCheck) bool {
	// Use reflection to determine whether the underlying type of the filtering attribute is a string or
	// array of strings. Mainly used because the property could be an SDK enum with underlying string type.
	val := reflect.ValueOf(target)
	valType := val.Type()

	for _, fVal := range filters {
		switch valType.Kind() {
		case reflect.Bool:
			fBool, err := strconv.ParseBool(fVal.(string))
			if err != nil {
				log.Println("[WARN] Filtering against Type Bool field with un-parsable string boolean form")
				return false
			}
			if val.Bool() == fBool {
				return true
			}
		case reflect.Int, reflect.Int32, reflect.Int64:
			// the target field is of type int, but the filter values list element type is string, users can supply string
			// or int like `values = [300, "3600"]` but terraform will converts to string, so use ParseInt
			fInt, err := strconv.ParseInt(fVal.(string), 10, 64)
			if err != nil {
				log.Println("[WARN] Filtering against Type Int field with non-int filter value")
				return false
			}
			if val.Int() == fInt {
				return true
			}
		case reflect.Float64:
			// same comment as above for Ints
			fFloat, err := strconv.ParseFloat(fVal.(string), 64)
			if err != nil {
				log.Println("[WARN] Filtering against Type Float field with non-float filter value")
				return false
			}
			if val.Float() == fFloat {
				return true
			}
		case reflect.String:
			if stringsEqual(val.String(), fVal.(string)) {
				return true
			}
		case reflect.Slice, reflect.Array:
			if valType.Elem().Kind() == reflect.String {
				arrLen := val.Len()
				for i := 0; i < arrLen; i++ {
					if stringsEqual(val.Index(i).String(), fVal.(string)) {
						return true
					}
				}
			}
		}
	}
	return false
}

func checkAndConvertMap(element interface{}) (map[string]interface{}, bool) {
	if tempWorkingMap, isOk := element.(map[string]interface{}); isOk {
		return tempWorkingMap, true
	}

	if stringToStrinMap, isOk := element.(map[string]string); isOk {
		return convertToObjectMap(stringToStrinMap), true
	}

	return nil, false
}

func convertToObjectMap(stringTostring map[string]string) map[string]interface{} {
	convertedMap := make(map[string]interface{}, len(stringTostring))
	for key, value := range stringTostring {
		convertedMap[key] = value
	}

	return convertedMap
}

func checkAndConvertNestedStructure(element interface{}) (map[string]interface{}, bool) {
	if convertedList, convertedListOk := element.([]interface{}); convertedListOk && len(convertedList) == 1 {
		workingMap, isOk := convertedList[0].(map[string]interface{})
		return workingMap, isOk
	}

	return nil, false
}

func getValueFromPath(item HclKeyValueObject, path []string) (targetVal interface{}, targetValOk bool) {
	workingMap := item
	tempWorkingMap := item
	var conversionOk bool
	for _, pathElement := range path[:len(path)-1] {
		// Defensive check for non existent values
		if workingMap[pathElement] == nil {
			return nil, false
		}
		// Check if it is map
		if tempWorkingMap, conversionOk = checkAndConvertMap(workingMap[pathElement]); !conversionOk {
			// if not map then it has to be a nested structure which is modeled as list with exactly one element of type map[string]interface{}
			if tempWorkingMap, conversionOk = checkAndConvertNestedStructure(workingMap[pathElement]); !conversionOk {
				return nil, false
			}
		}
		workingMap = tempWorkingMap
	}

	targetVal, targetValOk = workingMap[path[len(path)-1]]
	return
}
