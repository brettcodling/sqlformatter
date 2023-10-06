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
	var segments []tokens.Token
	tab := "\t"

	origSegments := tokens.Tokenize(query)

	for i, t := range origSegments {
		if t.TokenType != 0 {
			t.Index = i
			segments = append(segments, t)
		}
	}

	for i, t := range segments {
		highlighted := t.TokenValue
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

		if t.TokenType == tokens.TOKEN_TYPE_COMMENT || t.TokenType == tokens.TOKEN_TYPE_BLOCK_COMMENT {
			if t.TokenType == tokens.TOKEN_TYPE_BLOCK_COMMENT {
				indent := strings.Repeat(tab, indentLevel)
				ret += "\n" + indent
				highlighted = strings.ReplaceAll(highlighted, "\n", "\n"+indent)
			}

			ret += highlighted
			newline = true
			continue
		}

		if inlineParanthesis {
			if t.TokenValue == ")" {
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

			if t.TokenValue == "," {
				if inlineCount >= 30 {
					inlineCount = 0
					newline = true
				}
			}

			inlineCount += len(t.TokenValue)
		}

		if t.TokenValue == "(" {
			length := 0
			for j := 1; j <= 250; j++ {
				if len(segments) <= i+j {
					break
				}
				next := segments[i+j]

				if next.TokenValue == ")" {
					inlineParanthesis = true
					inlineCount = 0
					inlineIndented = false
					break
				}

				if next.TokenValue == ";" || next.TokenValue == "(" {
					break
				}

				if next.TokenType == tokens.TOKEN_TYPE_RESERVED_TOPLEVEL || next.TokenType == tokens.TOKEN_TYPE_RESERVED_NEWLINE || next.TokenType == tokens.TOKEN_TYPE_COMMENT || next.TokenType == tokens.TOKEN_TYPE_BLOCK_COMMENT {
					break
				}

				length += len(next.TokenValue)
			}

			if inlineParanthesis && length > 30 {
				incBlockIndent = true
				inlineIndented = true
				newline = true
			}

			if !inlineParanthesis {
				incBlockIndent = true
				newline = true
			}
		} else if t.TokenValue == ")" {
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
		} else if t.TokenType == tokens.TOKEN_TYPE_RESERVED_TOPLEVEL {
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

			if strings.Contains(t.TokenValue, " ") || strings.Contains(t.TokenValue, "\n") || strings.Contains(t.TokenValue, "\t") {
				re := regexp.MustCompile("\\s+")
				highlighted = re.ReplaceAllString(highlighted, " ")
			}

			if t.TokenValue == "LIMIT" && !inlineParanthesis {
				clauseLimit = true
			}
		} else if clauseLimit && t.TokenValue != "," && t.TokenType != tokens.TOKEN_TYPE_NUMBER && t.TokenType != tokens.TOKEN_TYPE_WHITESPACE {
			clauseLimit = false
		} else if t.TokenValue == "," && !inlineParanthesis {
			if clauseLimit {
				newline = false
				clauseLimit = false
			} else {
				newline = true
			}
		} else if t.TokenType == tokens.TOKEN_TYPE_RESERVED_NEWLINE {
			if !addedNewline {
				ret += "\n" + strings.Repeat(tab, indentLevel)
			}

			if strings.Contains(t.TokenValue, " ") || strings.Contains(t.TokenValue, "\n") || strings.Contains(t.TokenValue, "\t") {
				re := regexp.MustCompile("\\s+")
				highlighted = re.ReplaceAllString(highlighted, " ")
			}
		}

		if t.TokenValue == "." || t.TokenValue == "," || t.TokenValue == ";" {
			ret = strings.TrimRight(ret, " ")
		}

		ret += highlighted
		lastChar := ret[len(ret)-1:]
		if lastChar != " " && lastChar != "." {
			ret += " "
		}

		if lastChar == "(" {
			inlineParanthesis = true
		}

		if t.TokenValue == "(" || t.TokenValue == "." {
			ret = strings.TrimRight(ret, " ")
		}

		if t.TokenValue == "-" && len(segments) >= i && segments[i+1].TokenType == tokens.TOKEN_TYPE_NUMBER && i-1 >= 0 {
			prev := segments[i-1].TokenType
			if prev != tokens.TOKEN_TYPE_QUOTE &&
				prev != tokens.TOKEN_TYPE_BACKTICK_QUOTE &&
				prev != tokens.TOKEN_TYPE_WORD &&
				prev != tokens.TOKEN_TYPE_NUMBER {
				ret = strings.TrimRight(ret, " ")
			}
		}
	}

	return ret
}
