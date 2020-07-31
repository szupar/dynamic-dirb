package services

import "dynamic-dirb/internal/helper"

var parameters *helper.ParamValidator

func GetParameters() *helper.ParamValidator {
	return parameters
}

func SetParameters(newParameters *helper.ParamValidator) {
	parameters = newParameters
}
