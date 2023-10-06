package sqlformatter

import (
	"regexp"
	"strings"
	
	"github.com/brettcodling/sqlformatter/pkg/tokens"
)

func Format(query string) (formattedQuery string) {
	var ret string
	var indentLevel, inlineCount int
	var newline, inlineParanthesis, incSpecIndent, incBlockIndent, addedNewline, inlineIndented, clauseLimit bool
	var indentTypes []string
	var segments []tokens.token
	tab := "\t"

	origSegments := tokens.tokenize(query)

	for i, t := range origSegments {
		if t.tokenType != 0 {
			t.index = i
			segments = append(segments, t)
		}
	}

	for i, t := range segments {
		highlighted := t.tokenValue
		if incSpecIndent {
			indentLevel++
			incSpecIndent = false
			indentTypes = append([]string{"special"}, indentTypes...)
		}
		if incBlockIndent {
			indentLevel++
			incBlockIndent = false
			indentTypes = append([]string{"block"}, indentTypes...)
		}
		if newline {
			ret += "\n" + strings.Repeat(tab, indentLevel)
			newline = false
			addedNewline = true
		} else {
			addedNewline = false
		}

		if t.tokenType == 8 || t.tokenType == 9 {
			if t.tokenType == 9 {
				indent := strings.Repeat(tab, indentLevel)
				ret += "\n" + indent
				highlighted = strings.ReplaceAll(highlighted, "\n", "\n"+indent)
			}

			ret += highlighted
			newline = true
			continue
		}

		if inlineParanthesis {
			if t.tokenValue == ")" {
				ret = strings.TrimRight(ret, " ")

				if inlineIndented {
					indentTypes = indentTypes[1:]
					indentLevel--
					ret += "\n" + strings.Repeat(tab, indentLevel)
				}

				inlineParanthesis = false
				ret += highlighted + " "
				continue
			}

			if t.tokenValue == "," {
				if inlineCount >= 30 {
					inlineCount = 0
					newline = true
				}
			}

			inlineCount += len(t.tokenValue)
		}

		if t.tokenValue == "(" {
			length := 0
			for j := 1; j <= 250; j++ {
				if len(segments) <= i+j {
					break
				}
				next := segments[i+j]

				if next.tokenValue == ")" {
					inlineParanthesis = true
					inlineCount = 0
					inlineIndented = false
					break
				}

				if next.tokenValue == ";" || next.tokenValue == "(" {
					break
				}

				if next.tokenType == 5 || next.tokenType == 6 || next.tokenType == 8 || next.tokenType == 9 {
					break
				}

				length += len(next.tokenValue)
			}

			if inlineParanthesis && length > 30 {
				incBlockIndent = true
				inlineIndented = true
				newline = true
			}

			if t.index-1 >= 0 && t.index-1 < len(origSegments) && origSegments[t.index-1].tokenType != 0 {
				ret = strings.TrimRight(ret, " ")
			}

			if !inlineParanthesis {
				incBlockIndent = true
				newline = true
			}
		} else if t.tokenValue == ")" {
			ret = strings.TrimRight(ret, " ")

			indentLevel--

			for _, indentType := range indentTypes {
				indentTypes = indentTypes[1:]
				if indentType == "special" {
					indentLevel--
				} else {
					break
				}
			}

			if indentLevel < 0 {
				indentLevel = 0
			}

			if !addedNewline {
				ret += "\n" + strings.Repeat(tab, indentLevel)
			}
		} else if t.tokenType == 5 {
			incSpecIndent = true

			if len(indentTypes) > 0 && indentTypes[0] == "special" {
				indentLevel--
				indentTypes = indentTypes[1:]
			}

			newline = true
			if !addedNewline {
				ret += "\n" + strings.Repeat(tab, indentLevel)
			} else {
				ret = strings.TrimRight(ret, tab) + strings.Repeat(tab, indentLevel)
			}

			if strings.Contains(t.tokenValue, " ") || strings.Contains(t.tokenValue, "\n") || strings.Contains(t.tokenValue, "\t") {
				re := regexp.MustCompile("\\s+")
				highlighted = re.ReplaceAllString(highlighted, " ")
			}

			if t.tokenValue == "LIMIT" && !inlineParanthesis {
				clauseLimit = true
			}
		} else if clauseLimit && t.tokenValue != "," && t.tokenType != 10 && t.tokenType != 0 {
			clauseLimit = false
		} else if t.tokenValue == "," && !inlineParanthesis {
			if clauseLimit {
				newline = false
				clauseLimit = false
			} else {
				newline = true
			}
		} else if t.tokenType == 6 {
			if !addedNewline {
				ret += "\n" + strings.Repeat(tab, indentLevel)
			}

			if strings.Contains(t.tokenValue, " ") || strings.Contains(t.tokenValue, "\n") || strings.Contains(t.tokenValue, "\t") {
				re := regexp.MustCompile("\\s+")
				highlighted = re.ReplaceAllString(highlighted, " ")
			}
		} else if t.tokenType == 7 {
			if i-1 >= 0 && i-1 < len(segments) && segments[i-1].tokenType != 7 {
				if t.index-1 >= 0 && t.index-1 < len(origSegments) && origSegments[t.index-1].tokenType != 0 {
					ret = strings.TrimRight(ret, " ")
				}
			}
		}

		if t.tokenValue == "." || t.tokenValue == "," || t.tokenValue == ";" {
			ret = strings.TrimRight(ret, " ")
		}

		ret += highlighted + " "

		if t.tokenValue == "(" || t.tokenValue == "." {
			ret = strings.TrimRight(ret, " ")
		}

		if t.tokenValue == "-" && len(segments) >= i && segments[i+1].tokenType == 10 && i-1 >= 0 {
			prev := segments[i-1].tokenType
			if prev != 2 && prev != 3 && prev != 1 && prev != 10 {
				ret = strings.TrimRight(ret, " ")
			}
		}
	}

	ret = strings.TrimSpace(strings.ReplaceAll(ret, tab, "  "))

	return ret
}
