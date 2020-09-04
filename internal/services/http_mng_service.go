package services

func IsResourceExist(validCode map[int]bool, codeToValidate int) bool {

	return validCode[codeToValidate]
}
