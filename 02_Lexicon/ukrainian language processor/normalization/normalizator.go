package normalization

import (
	"fmt"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"unicode"
)

const blackList = ".,-~?!\"'`;:()<>[]{}\\|/=_+*&^%$#@"

type Normalizator struct {
	t transform.Transformer
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

func New() *Normalizator {
	return &Normalizator{t:transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)}
}

func (n *Normalizator) Normalize(word string) string {
	for _, r := range blackList {
		word = strings.Replace(word, fmt.Sprintf("%c", r), "", -1)
	}
	res, _, _ := transform.String(n.t, word)
	return strings.ToLower(res)
}
