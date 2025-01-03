package reflect

import (
	"reflect"
	"strings"
)

// GetNonZeroFields 获取结构体中非零值字段
func GetNonZeroFields(data interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	v := reflect.ValueOf(data)

	// 如果是指针，获取其指向的值
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 确保是结构体
	if v.Kind() != reflect.Struct {
		return result
	}

	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i)

		// 跳过空值和特定字段
		if value.IsZero() || field.Name == "Did" {
			continue
		}

		// 获取 json tag
		tag := field.Tag.Get("json")
		if tag == "" {
			tag = strings.ToLower(field.Name)
		}

		result[tag] = value.Interface()
	}

	return result
}
