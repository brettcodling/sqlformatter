package sqlformatter

import (
	"regexp"
	"strings"

	"github.com/brettcodling/sqlformatter/pkg/tokens"
)

const (
	Indent_type_special = iota
	Indent_type_block
)

type formattedQuery struct {
	output string
}

func (f *formattedQuery) addNewLine(indent int) {
	f.output += "\n" + strings.Repeat("\t", indent)
}

func (f *formattedQuery) append(s string) {
	f.output += s
}

func (f *formattedQuery) trimRight() {
	f.output = strings.TrimRight(f.output, " ")
}

func Format(query string) string {
	var formatted formattedQuery
	var indentLevel, inlineCount int
	var newline, inlineParanthesis, incSpecIndent, incBlockIndent, addedNewline, inlineIndented, clauseLimit bool
	var indentTypes []int
	var segments []tokens.Token

	origSegments := tokens.Tokenize(tokens.Query(query))

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
			indentTypes = append([]int{Indent_type_special}, indentTypes...)
		}
		if incBlockIndent {
			indentLevel++
			incBlockIndent = false
			indentTypes = append([]int{Indent_type_block}, indentTypes...)
		}
		if newline {
			formatted.addNewLine(indentLevel)
			newline = false
			addedNewline = true
		} else {
			addedNewline = false
		}

		if t.TokenType == tokens.Token_type_comment || t.TokenType == tokens.Token_type_block_comment {
			if t.TokenType == tokens.Token_type_block_comment {
				indent := strings.Repeat("\t", indentLevel)
				formatted.addNewLine(indentLevel)
				highlighted = strings.ReplaceAll(highlighted, "\n", "\n"+indent)
			}

			formatted.append(highlighted)
			newline = true
			continue
		}

		if inlineParanthesis {
			if t.TokenValue == ")" {
				formatted.trimRight()

				if inlineIndented {
					indentTypes = indentTypes[1:]
					indentLevel--
					formatted.addNewLine(indentLevel)
				}

				inlineParanthesis = false
				formatted.append(highlighted + " ")
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

				if next.TokenType == tokens.Token_type_reserved_toplevel || next.TokenType == tokens.Token_type_reserved_newline || next.TokenType == tokens.Token_type_comment || next.TokenType == tokens.Token_type_block_comment {
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
			formatted.trimRight()

			indentLevel--

			for _, indentType := range indentTypes {
				indentTypes = indentTypes[1:]
				if indentType == Indent_type_special {
					indentLevel--
				} else {
					break
				}
			}

			if indentLevel < 0 {
				indentLevel = 0
			}

			if !addedNewline {
				formatted.addNewLine(indentLevel)
			}
		} else if t.TokenType == tokens.Token_type_reserved_toplevel {
			incSpecIndent = true

			if len(indentTypes) > 0 && indentTypes[0] == Indent_type_special {
				indentLevel--
				indentTypes = indentTypes[1:]
			}

			newline = true
			if !addedNewline {
				formatted.addNewLine(indentLevel)
			} else {
				formatted.trimRight()
				formatted.append(strings.Repeat("\t", indentLevel))
			}

			if strings.Contains(t.TokenValue, " ") || strings.Contains(t.TokenValue, "\n") || strings.Contains(t.TokenValue, "\t") {
				re := regexp.MustCompile("\\s+")
				highlighted = re.ReplaceAllString(highlighted, " ")
			}

			if t.TokenValue == "LIMIT" && !inlineParanthesis {
				clauseLimit = true
			}
		} else if clauseLimit && t.TokenValue != "," && t.TokenType != tokens.Token_type_number && t.TokenType != tokens.Token_type_whitespace {
			clauseLimit = false
		} else if t.TokenValue == "," && !inlineParanthesis {
			if clauseLimit {
				newline = false
				clauseLimit = false
			} else {
				newline = true
			}
		} else if t.TokenType == tokens.Token_type_reserved_newline {
			if !addedNewline {
				formatted.addNewLine(indentLevel)
			}

			if strings.Contains(t.TokenValue, " ") || strings.Contains(t.TokenValue, "\n") || strings.Contains(t.TokenValue, "\t") {
				re := regexp.MustCompile("\\s+")
				highlighted = re.ReplaceAllString(highlighted, " ")
			}
		}

		if t.TokenValue == "." || t.TokenValue == "," || t.TokenValue == ";" {
			formatted.trimRight()
		}

		formatted.append(highlighted)
		lastChar := formatted.output[len(formatted.output)-1:]
		if lastChar != " " && lastChar != "." {
			formatted.append(" ")
		}

		if lastChar == "(" {
			inlineParanthesis = true
		}

		if t.TokenValue == "(" || t.TokenValue == "." {
			formatted.trimRight()
		}

		if t.TokenValue == "-" && len(segments) >= i && segments[i+1].TokenType == tokens.Token_type_number && i-1 >= 0 {
			prev := segments[i-1].TokenType
			if prev != tokens.Token_type_quote &&
				prev != tokens.Token_type_backtick_quote &&
				prev != tokens.Token_type_word &&
				prev != tokens.Token_type_number {
				formatted.trimRight()
			}
		}
	}

	return formatted.output
}
