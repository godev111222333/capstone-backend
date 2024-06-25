package api

import (
	"github.com/go-playground/validator/v10"
)

var CustomValidations = map[string]validator.Func{
	"id_card":         validationFuncIDCard,
	"phone_number":    validationFuncPhoneNumber,
	"driving_license": validationFuncDrivingLicense,
	"license_plate":   validationFuncLicensePlate,
}

var validationFuncIDCard validator.Func = func(fl validator.FieldLevel) bool {
	idCard, ok := fl.Field().Interface().(string)
	if ok {
		l := len(idCard)
		if l == 0 || (l >= 9 && l <= 12) {
			return true
		}

		return false
	}

	return true
}

var validationFuncPhoneNumber validator.Func = func(fl validator.FieldLevel) bool {
	phoneNumber, ok := fl.Field().Interface().(string)
	if ok {
		l := len(phoneNumber)
		if l == 0 || l == 10 {
			return true
		}

		return false
	}

	return true
}

var validationFuncDrivingLicense validator.Func = func(fl validator.FieldLevel) bool {
	drivingLicense, ok := fl.Field().Interface().(string)
	if ok {
		l := len(drivingLicense)
		if l == 0 || l == 12 {
			return true
		}

		return false
	}

	return true
}

var validationFuncLicensePlate validator.Func = func(fl validator.FieldLevel) bool {
	drivingLicense, ok := fl.Field().Interface().(string)
	if ok {
		l := len(drivingLicense)
		if l == 0 || l == 7 || l == 8 {
			return true
		}

		return false
	}

	return true
}
