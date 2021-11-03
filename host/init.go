package host

import (
	"github.com/youthlin/t"
)

var ts = t.Global()

// ResetTranslations
//  t.Bind(domain, path)
//  ResetTranslations(t.D(domain))
func ResetTranslations(newTs *t.Translations) {
	ts = newTs
}
