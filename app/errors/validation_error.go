package errors

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func BindAndValidate(c *gin.Context, obj any) []*ValidationError {
	err := c.ShouldBindJSON(obj)
	if err == nil {
		return nil
	}

	var validationErrors []*ValidationError

	if unmarshalErrors, ok := err.(*json.UnmarshalTypeError); ok {
		validationErrors = append(validationErrors, &ValidationError{
			Field:   unmarshalErrors.Field,
			Message: fmt.Sprintf("Invalid data type: expected '%v' but got '%v'", unmarshalErrors.Type, unmarshalErrors.Value),
		})
	}

	if syntaxError, ok := err.(*json.SyntaxError); ok {
		validationErrors = append(validationErrors, &ValidationError{
			Message: fmt.Sprintf("Syntax error: %v", syntaxError.Error()),
		})
	}

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			validationErrors = append(validationErrors, &ValidationError{
				Field:   getFieldName(fe, obj),
				Message: getErrorMsg(fe),
			})
		}
	}

	if len(validationErrors) == 0 {
		validationErrors = append(validationErrors, &ValidationError{
			Message: "Missing request body",
		})
	}

	return validationErrors
}

func getFieldName(fe validator.FieldError, obj any) string {
	fieldName := fe.Field()

	field, _ := reflect.TypeOf(obj).Elem().FieldByName(fieldName)
	jsonName := field.Tag.Get("json")
	if jsonName != "" {
		fieldName = jsonName
	}
	return fieldName
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "lte", "max":
		return "Should be less than or equal to " + fe.Param()
	case "gte", "min":
		return "Should be greater than or equal to " + fe.Param()
	case "lt":
		return "Should be less than " + fe.Param()
	case "gt":
		return "Should be greater than " + fe.Param()
	case "eq":
		return "Should be equal to " + fe.Param()
	case "ne":
		return "Should be not equal to " + fe.Param()
	case "email":
		return "Should be a valid email address"
	case "oneof":
		return "Should be one of " + strings.Join(strings.Split(fe.Param(), " "), ", ")
	case "date":
		return "Should be a valid date with format YYYY-MM-DD"
	case "beforeToday":
		return "Should be a valid date before today"
	case "phoneNumber":
		return "Should be a valid phone number"
	}
	return fe.Error()
}
