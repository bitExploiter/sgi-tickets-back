package toolbox

import (
	"reflect"

	"github.com/go-playground/validator"
)

func FormatValidationErrors(s interface{}) (map[string]string, error) {
	validate := validator.New()
	err := validate.Struct(s)
	if err != nil {
		errorsMap := make(map[string]string)
		t := reflect.TypeOf(s)
		for _, err := range err.(validator.ValidationErrors) {
			field, _ := t.FieldByName(err.StructField())
			jsonTag := field.Tag.Get("json")
			if jsonTag == "" {
				jsonTag = err.Field()
			}
			message := err.Tag()
			errorsMap[jsonTag] = message
		}
		return errorsMap, err
	}
	return nil, nil
}
