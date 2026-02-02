package server

import (
	// "github.com/go-playground/locales/currency"

	//TODO:: will use the go-currency package later

	"github.com/go-playground/validator/v10"
	"{{ cookiecutter.module_path }}/utils"
)

var validCurrency validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		//  check currency is supported
		return utils.IsSupportedCurrency(currency)
	}
	return false
}

// type CustomValidator struct {
// 	Validator *validator.Validate
// }

// var storeRef *db.Store

// func InitValidator(store *db.Store) {
// 	storeRef = store

// 	v := validator.New()

// 	// register custom validation functions
// 	v.RegisterValidation("validCategories", validateCategories)
// 	v.RegisterValidation("validAmenities", validateAmenities)
// 	v.RegisterValidation("validSafetyItems", validateSafetyItems)

// 	// Tell Gin to use OUR validator instead of default!
// 	binding.Validator = &CustomValidator{Validator: v}
// }

// func (cv *CustomValidator) ValidateStruct(obj interface{}) error {
// 	if cv.Validator == nil {
// 		return nil
// 	}
// 	return cv.Validator.Struct(obj)
// }

// func (cv *CustomValidator) Engine() interface{} {
// 	return cv.Validator
// }

// func validateCategories(fl validator.FieldLevel) bool {
// 	fields := fl.Field()

// 	for i := 0; i < fields.Len(); i++ {
// 		value := fields.Index(i).String()

// 		_, err := (*storeRef).GetCategoryName(context.Background(), value)
// 		if err != nil {
// 			return false
// 		}
// 	}

// 	return true
// }

// func validateAmenities(fl validator.FieldLevel) bool {
// 	fields := fl.Field()

// 	for i := 0; i < fields.Len(); i++ {
// 		value := fields.Index(i).String()

// 		_, err := (*storeRef).GetAmenityName(context.Background(), value)
// 		if err != nil {
// 			return false
// 		}
// 	}
// 	return true
// }

// func validateSafetyItems(fl validator.FieldLevel) bool {
// 	fields := fl.Field()

// 	for i := 0; i < fields.Len(); i++ {
// 		value := fields.Index(i).String()

// 		_, err := (*storeRef).GetSafetyItemName(context.Background(), value)
// 		if err != nil {
// 			return false
// 		}
// 	}
// 	return true
// }
