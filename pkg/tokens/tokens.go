package tokens

import (
	"regexp"
	"strings"
)

type token struct {
	tokenValue string
	tokenType  int
	index      int
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
	re = regexp.MustCompile("^([0-9]+(\\.[0-9]+)?|0x[0-9a-fA-F]+|0b[01]+)($|\\s|\"'\\x60|" + regexBoundaries + ")")
	if match := re.FindString(query); match != "" {
		return token{
			tokenValue: match,
			tokenType:  10,
		}
	}

	// Boundary
	re = regexp.MustCompile("^" + regexBoundaries)
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
			if cachedToken, exists := tokenCache[cacheKey]; exists {
				t = cachedToken
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
