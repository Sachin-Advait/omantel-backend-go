package objecthelper

import "unicode"

func IsStrictlyAlphNumeric(word string) bool {

	isNumber := false
	isAplpha := false

	for _, v := range word {
		if unicode.IsNumber(v) {
			isNumber = true
			continue
		}

		if unicode.IsLetter(v) {
			isAplpha = true
			continue
		}

		return false
	}

	if isAplpha && isNumber {
		return true
	}

	return false
}

func ContainsAlphNumeric(word string) bool {

	isNumber := false
	isAplpha := false

	for _, v := range word {
		if unicode.IsNumber(v) {
			isNumber = true
		}

		if unicode.IsLetter(v) {
			isAplpha = true
		}
	}

	if isAplpha && isNumber {
		return true
	}

	return false
}

func IsStrictlyNumeric(num string) bool {

	for _, v := range num {
		if !unicode.IsNumber(v) {
			return false
		}
	}

	return true
}

func ContainsNumeric(num string) bool {

	for _, v := range num {
		if unicode.IsNumber(v) {
			return true
		}
	}

	return false
}

func IsStrictlyAlpha(alpha string) bool {

	for _, v := range alpha {
		if !unicode.IsLetter(v) {
			return false
		}
	}

	return true
}

func ContainsAlpha(alpha string) bool {

	for _, v := range alpha {
		if unicode.IsLetter(v) {
			return true
		}
	}

	return false
}
