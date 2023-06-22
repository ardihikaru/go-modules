package common

import (
	"strings"
	"unicode/utf8"
)

func SanitizePhone(phone string, withPlusSymbol *bool) string {
	// removes `+` symbol if validation enabled
	if withPlusSymbol != nil && !(*withPlusSymbol) && phone[0:1] == "+" {
		_, i := utf8.DecodeRuneInString(phone)
		phone = phone[i:]
	}
	// replace minus symbol and spaces
	phone = strings.Replace(phone, "-", "", -1)
	phone = strings.Replace(phone, " ", "", -1)

	return phone
}
