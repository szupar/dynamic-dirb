package services

func IsResourceExist(validCode map[int]bool, codeToValidate int) bool {

	if validCode[codeToValidate] {
		return true
	}

	return false
}
