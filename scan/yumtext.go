package scan

import (
	"regexp"
	"strings"
)

type YUM_INFO_TYPE int

const (
	YUM_UNKNOWN YUM_INFO_TYPE = 0

	YUM_NAME    YUM_INFO_TYPE = 1
	YUM_VERSION YUM_INFO_TYPE = 2
)

var yumNameExp *regexp.Regexp
var yumVersionExp *regexp.Regexp
var yumVersionExp2 *regexp.Regexp
var yumTextSplitExp *regexp.Regexp

func init() {
	yumNameExp = regexp.MustCompile(`([\w\-\+]+)\.(i[3-6]86|noarch|x86|x86_64|athlon)$`)
	yumVersionExp = regexp.MustCompile(`([:\w\.\-]+)\.(rh)?el[567][_\.\w]*`)
	yumVersionExp2 = regexp.MustCompile(`[\d\-\.+:]`)
	yumTextSplitExp = regexp.MustCompile(`^([\w\.\-\+]+)\s+([:\w\.\-]+)\s+(\S*)$`)
}

func parseYumText(txt string) YUM_INFO_TYPE {
	m := yumNameExp.FindAllStringSubmatch(txt, 1)
	//fmt.Println("match:", m)

	if yumNameExp.MatchString(txt) && len(m) > 0 {
		var score int
		for idx, c := range m[0][1] {
			if (c >= 'a' && c <= 'z') ||
				(c >= 'A' && c <= 'Z') || c == '_' || c == '-' {
				if idx == 0 {
					score += 5
				} else {
					score += 1
				}
			} else if c >= '0' && c <= '9' {
				score -= 2
			} else if c == '.' {
				score -= 5
			}
		}
		if score > 0 {
			return YUM_NAME
		} else {
			//			fmt.Println("TESTA", m[0][1], "score:", score)
		}
	}

	m1 := yumVersionExp.FindAllStringSubmatch(txt, 1)
	if yumVersionExp.MatchString(txt) && len(m1) > 0 {
		var score int
		var numbercount int
		for idx, c := range m1[0][1] {
			if (c >= 'a' && c <= 'z') ||
				(c >= 'A' && c <= 'Z') || c == '_' {
				score -= 2
			} else if c >= '0' && c <= '9' {
				if idx == 0 {
					score += 5
				} else {
					score += 1
				}
			} else if c == '.' {
				score += 2
			}
		}
		if len(m1[0]) > 2 {
			score++
		}
		if score > 0 && numbercount > 0 {
			return YUM_VERSION
		} else {
			//			fmt.Println("TESTB", m1[0][1], "score:", score)
		}
	}
	if yumVersionExp2.MatchString(txt) {
		return YUM_VERSION
	}

	return YUM_UNKNOWN
}

type yumvalue struct {
	txt string
	typ YUM_INFO_TYPE
}

func preProcessText(line string) []yumvalue {
	var r []yumvalue
	line = strings.TrimSpace(line)
	m := yumTextSplitExp.FindAllStringSubmatch(line, 1)
	if len(m) > 0 {
		for k, s := range m[0] {
			if k == 0 {
				continue
			}
			r = append(r, yumvalue{
				txt: s,
				typ: parseYumText(s),
			})
		}
		return r
	}
	return r
}
