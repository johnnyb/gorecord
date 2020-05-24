package inflect

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func init() {
	RegisterSingularPlural("person", "people")
	RegisterSingularPlural("sheep", "sheep")
	RegisterNonTitleized("a")
	RegisterNonTitleized("the")
}

var Plurals map[string]string = map[string]string{}
var Singulars map[string]string = map[string]string{}
var NonTitleized map[string]bool = map[string]bool{}

func RegisterNonTitleized(str string) {
	NonTitleized[str] = true
}

func RegisterSingularPlural(singular string, plural string) {
	Plurals[singular] = plural
	Singulars[plural] = singular
}

func Titleize(str string) string {
	comps := Componentize(str)
	newcomps := []string{}
	for _, comp := range comps {
		if NonTitleized[comp] {
			newcomps = append(newcomps, comp)
		} else {
			newcomps = append(newcomps, CapitalizeFirstLetter(comp))
		}
	}
	return strings.Join(newcomps, " ")
}

func Camelize(str string) string {
	return str

}

func Humanize(str string) string {
	return str

}

func Tableize(str string) string {
	return Pluralize(Underscore(str))
}

func IsSingular(str string) bool {
	last := LastComponent(str)
	if last == "" {
		return false
	}
	custom := Plurals[last]
	if custom != "" {
		return true
	}
	custom = Singulars[last]
	if custom != "" {
		return false
	}

	_, lr := LastRune(last)
	return lr != 's'
}

func IsPlural(str string) bool {
	last := LastComponent(str)
	if str == "" {
		return false
	}
	custom := Singulars[str]
	if custom != "" {
		return true
	}
	custom = Plurals[str]
	if custom != "" {
		return false
	}
	_, lr := LastRune(last)
	return lr == 's'
}

func Pluralize(str string) string {
	last := LastComponent(str)
	if IsPlural(last) {
		return str
	}
	custom := Plurals[last]

	if custom != "" {
		if IsLastComponentCapitalized(str, last) {
			custom = CapitalizeFirstLetter(custom)
		}
		newstr := str[0:(len(str) - len(last))] + custom
		return newstr
	}

	return str + "s"
}

func Componentize(str string) []string {
	curstr := ""
	results := []string{}

	for _, rn := range str {
		if rn == '-' || rn == '_' || unicode.IsSpace(rn) {
			if curstr != "" {
				results = append(results, curstr)
				curstr = ""
			}
		} else if unicode.IsUpper(rn) {
			if curstr != "" {
				results = append(results, curstr)
				curstr = ""
			}
			curstr = curstr + string(unicode.ToLower(rn))
		} else {
			curstr = curstr + string(rn)
		}
	}
	if curstr != "" {
		results = append(results, curstr)
	}

	return results
}

func CapitalizeFirstLetter(str string) string {
	firstRune, sz := utf8.DecodeRuneInString(str)
	if unicode.IsUpper(firstRune) {
		return str
	}
	firstRune = unicode.ToUpper(firstRune)

	return string(firstRune) + str[sz:]
}

func LastRune(str string) (int, rune) {
	var idx int
	var ch rune
	for ridx, r := range str {
		idx = ridx
		ch = r
	}
	return idx, ch
}

func IsLastComponentCapitalized(str string, last string) bool {
	lastrune, _ := utf8.DecodeRuneInString(str[(len(str) - len(last)):])
	return unicode.IsUpper(lastrune)
}

func LastComponent(str string) string {
	comps := Componentize(str)
	if len(comps) == 0 {
		return ""
	}
	return comps[len(comps) - 1]
}

func Singularize(str string) string {
	last := LastComponent(str)
	if IsSingular(last) {
		return str
	}
	custom := Singulars[last]

	if custom != "" {
		if IsLastComponentCapitalized(str, last) {
			custom = CapitalizeFirstLetter(custom)
		}
		newstr := str[0:(len(str) - len(last))] + custom
		return newstr
	}

	return str[0:len(str) - 1]
}

func Underscore(str string) string {
	return strings.Join(Componentize(str), "_")
}

