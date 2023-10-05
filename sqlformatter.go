package sqlformatter

import (
	"regexp"
	"strings"
)

func Format(query string) (formattedQuery string) {
	var ret string
	var indentLevel, inlineCount int
	var newline, inlineParanthesis, incSpecIndent, incBlockIndent, addedNewline, inlineIndented, clauseLimit bool
	var indentTypes []string
	var tokens []token
	tab := "\t"

	origTokens := tokenize(query)

	for i, t := range origTokens {
		if t.tokenType != 0 {
			t.index = i
			tokens = append(tokens, t)
		}
	}

	for i, t := range tokens {
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
				if len(tokens) <= i+j {
					break
				}
				next := tokens[i+j]

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

			if t.index-1 >= 0 && t.index-1 < len(origTokens) && origTokens[t.index-1].tokenType != 0 {
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
			if i-1 >= 0 && i-1 < len(tokens) && tokens[i-1].tokenType != 7 {
				if t.index-1 >= 0 && t.index-1 < len(origTokens) && origTokens[t.index-1].tokenType != 0 {
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

		if t.tokenValue == "-" && len(tokens) >= i && tokens[i+1].tokenType == 10 && i-1 >= 0 {
			prev := tokens[i-1].tokenType
			if prev != 2 && prev != 3 && prev != 1 && prev != 10 {
				ret = strings.TrimRight(ret, " ")
			}
		}
	}

	ret = strings.TrimSpace(strings.ReplaceAll(ret, tab, "  "))

	return ret
}

func getNextToken(query string, t token) token {
	// Whitespace
	re := regexp.MustCompile("^\\s+")
	if match := re.FindString(query); match != "" {
		return token{
			tokenValue: match,
			tokenType:  0,
		}
	}

	// Comment
	if query[:1] == "#" || (len(query) > 1 && (query[:2] == "--" || query[:2] == "/*")) {
		var last, tokenType int
		if query[:1] == "#" || query[:1] == "-" {
			last = strings.Index(query, "\n")
			tokenType = 8
		} else {
			last = strings.Index(query[2:], "*/") + 2
			tokenType = 9
		}

		if last == -1 {
			last = len(query)
		}

		return token{
			tokenValue: query[:last],
			tokenType:  tokenType,
		}
	}

	// Quoted
	if query[:1] == "\"" || query[:1] == "'" || query[:1] == "`" || query[:1] == "[" {
		tokenType := 2
		if query[:1] == "`" || query[:1] == "[" {
			tokenType = 3
		}
		return token{
			tokenValue: getQuotedString(query),
			tokenType:  tokenType,
		}
	}

	// User defined
	if len(query) > 1 && (query[:1] == "@" || query[:1] == ":") {
		ret := token{
			tokenValue: "",
			tokenType:  12,
		}

		if query[1:2] == "\"" || query[1:2] == "'" || query[1:2] == "`" {
			ret.tokenValue = query[:1] + getQuotedString(query[1:])
		} else {
			re := regexp.MustCompile("^(" + query[:1] + "[a-zA-Z0-9\\._\\$]+)")
			ret.tokenValue = re.FindString(query)
		}

		if ret.tokenValue != "" {
			return ret
		}
	}

	// Number
	re = regexp.MustCompile("([0-9]+(\\.[0-9]+)?|0x[0-9a-fA-F]+|0b[01]+)($|\\s|\"'\\x60|" + regexBoundaries + ")")
	if match := re.FindString(query); match != "" {
		return token{
			tokenValue: match,
			tokenType:  10,
		}
	}

	// Boundary
	re = regexp.MustCompile(regexBoundaries)
	if match := re.FindString(query); match != "" {
		return token{
			tokenValue: match,
			tokenType:  7,
		}
	}

	upper := strings.ToUpper(query)
	// Reserved from prefixed by '.'
	if t.tokenValue == "" || t.tokenValue != "." {
		re := regexp.MustCompile("^" + regexReservedToplevel + "($|\\s|" + regexBoundaries + ")")
		if match := re.FindString(upper); match != "" {
			return token{
				tokenValue: query[:len(match)],
				tokenType:  5,
			}
		}
		re = regexp.MustCompile("^" + regexReservedNewline + "($|\\s|" + regexBoundaries + ")")
		if match := re.FindString(upper); match != "" {
			return token{
				tokenValue: query[:len(match)],
				tokenType:  6,
			}
		}
		re = regexp.MustCompile("^" + regexReserved + "($|\\s|" + regexBoundaries + ")")
		if match := re.FindString(upper); match != "" {
			return token{
				tokenValue: query[:len(match)],
				tokenType:  4,
			}
		}
	}

	// Function suffixed by '('
	re = regexp.MustCompile("^(" + regexReserved + "[(]|\\s|[)])")
	if match := re.FindString(upper); match != "" {
		return token{
			tokenValue: query[:len(match)-1],
			tokenType:  4,
		}
	}

	// Non reserved
	re = regexp.MustCompile("^(.*?)($|\\s|[\"'\\x60]|" + regexBoundaries + ")")
	return token{
		tokenValue: re.FindString(query),
		tokenType:  1,
	}
}

func getQuotedString(query string) (quoted string) {
	re := regexp.MustCompile("^(((\\x60[^\\x60]*($|\\x60))+)|((\\[[^\\]]*($|\\]))(\\][^\\]]*($|\\]))*)|((\"[^\"\\\\]*(?:\\\\.[^\"\\\\]*)*(\"|$))+)|(('[^'\\\\]*(?:\\\\.[^'\\\\]*)*('|$))+))")
	quoted = re.FindString(query)

	return
}

func initialize() {
	regexBoundaries = joinQuotedRegexp(boundaries)
	regexReserved = joinQuotedRegexp(reserved)
	regexReservedToplevel = strings.ReplaceAll(joinQuotedRegexp(reservedToplevel), " ", "\\s+")
	regexReservedNewline = strings.ReplaceAll(joinQuotedRegexp(reservedNewline), " ", "\\s+")
	regexFunction = joinQuotedRegexp(functions)
}

func joinQuotedRegexp(values []string) string {
	var quoted []string
	for _, value := range values {
		quoted = append(quoted, regexp.QuoteMeta(value))
	}

	return "(" + strings.Join(quoted, "|") + ")"
}

func tokenize(query string) (tokens []token) {
	initialize()

	origLength := len(query)
	oldLength := origLength + 1
	currLength := origLength

	var t token

	for {
		if currLength < 1 {
			break
		}

		if oldLength <= currLength {
			tokens = append(tokens, token{
				tokenValue: query,
				tokenType:  11,
			})

			return
		}
		oldLength = currLength

		var cacheKey string
		if currLength >= 15 {
			cacheKey = query[0:15]
		}

		var tokenLength int
		if cacheKey != "" {
			var exists bool
			t, exists = tokenCache[cacheKey]
			if exists {
				tokenLength = len(t.tokenValue)
			}
		}
		if tokenLength < 1 {
			t = getNextToken(query, t)
			tokenLength = len(t.tokenValue)

			if cacheKey != "" && tokenLength < 15 {
				if tokenCache == nil {
					tokenCache = make(map[string]token)
				}
				tokenCache[cacheKey] = t
			}
		}

		tokens = append(tokens, t)
		query = query[tokenLength:]
		currLength -= tokenLength
	}

	return
}
