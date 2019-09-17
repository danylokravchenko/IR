package detector

import (
	"github.com/kapsteur/franco"
)

// return language code
func DetectLanguage(s string) string {
	res := franco.DetectOne(s)
	return res.Code
}
