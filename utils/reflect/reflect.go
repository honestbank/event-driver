package reflect

import (
	"reflect"
)

func GetType(object interface{}) string {
	pointers := ""
	t := reflect.TypeOf(object)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
		pointers += "*"
	}

	return pointers + t.Name()
}
