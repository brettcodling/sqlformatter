package tokens

import (
	"regexp"
	"strings"
)

type Token struct {
	TokenValue string
	TokenType  int
	Index      int
}

func getNextToken(query string, t Token) Token {
	// Whitespace
	re := regexp.MustCompile("^\\s+")
	if match := re.FindString(query); match != "" {
		return Token{
			TokenValue: match,
			TokenType:  TOKEN_TYPE_WHITESPACE,
		}
	}

	// Comment
	if query[:1] == "#" || (len(query) > 1 && (query[:2] == "--" || query[:2] == "/*")) {
		var last, tokenType int
		if query[:1] == "#" || query[:1] == "-" {
			last = strings.Index(query, "\n")
			tokenType = TOKEN_TYPE_COMMENT
		} else {
			last = strings.Index(query[2:], "*/") + 2
			tokenType = TOKEN_TYPE_BLOCK_COMMENT
		}

		if last == -1 {
			last = len(query)
		}

		return Token{
			TokenValue: query[:last],
			TokenType:  tokenType,
		}
	}

	// Quoted
	if query[:1] == "\"" || query[:1] == "'" || query[:1] == "`" || query[:1] == "[" {
		tokenType := TOKEN_TYPE_QUOTE
		if query[:1] == "`" || query[:1] == "[" {
			tokenType = TOKEN_TYPE_BACKTICK_QUOTE
		}
		return Token{
			TokenValue: getQuotedString(query),
			TokenType:  tokenType,
		}
	}

	// User defined
	if len(query) > 1 && (query[:1] == "@" || query[:1] == ":") {
		ret := Token{
			TokenValue: "",
			TokenType:  TOKEN_TYPE_VARIABLE,
		}

		if query[1:2] == "\"" || query[1:2] == "'" || query[1:2] == "`" {
			ret.TokenValue = query[:1] + getQuotedString(query[1:])
		} else {
			re := regexp.MustCompile("^(" + query[:1] + "[a-zA-Z0-9\\._\\$]+)")
			ret.TokenValue = re.FindString(query)
		}

		if ret.TokenValue != "" {
			return ret
		}
	}

	// Number
	re = regexp.MustCompile("^([0-9]+(\\.[0-9]+)?|0x[0-9a-fA-F]+|0b[01]+)($|\\s|\"'\\x60|" + regexBoundaries + ")")
	if match := re.FindString(query); match != "" {
		return Token{
			TokenValue: match,
			TokenType:  TOKEN_TYPE_NUMBER,
		}
	}

	// Boundary
	re = regexp.MustCompile("^" + regexBoundaries)
	if match := re.FindString(query); match != "" {
		return Token{
			TokenValue: match,
			TokenType:  TOKEN_TYPE_BOUNDARY,
		}
	}

	upper := strings.ToUpper(query)
	// Reserved from prefixed by '.'
	if t.TokenValue == "" || t.TokenValue != "." {
		re := regexp.MustCompile("^" + regexReservedToplevel + "($|\\s|" + regexBoundaries + ")")
		if match := re.FindString(upper); match != "" {
			return Token{
				TokenValue: query[:len(match)],
				TokenType:  TOKEN_TYPE_RESERVED_TOPLEVEL,
			}
		}
		re = regexp.MustCompile("^" + regexReservedNewline + "($|\\s|" + regexBoundaries + ")")
		if match := re.FindString(upper); match != "" {
			return Token{
				TokenValue: query[:len(match)],
				TokenType:  TOKEN_TYPE_RESERVED_NEWLINE,
			}
		}
		re = regexp.MustCompile("^" + regexReserved + "($|\\s|" + regexBoundaries + ")")
		if match := re.FindString(upper); match != "" {
			return Token{
				TokenValue: query[:len(match)],
				TokenType:  TOKEN_TYPE_RESERVED,
			}
		}
	}

	// Function suffixed by '('
	re = regexp.MustCompile("^(" + regexReserved + "[(]|\\s|[)])")
	if match := re.FindString(upper); match != "" {
		return Token{
			TokenValue: query[:len(match)-1],
			TokenType:  TOKEN_TYPE_RESERVED,
		}
	}

	// Non reserved
	re = regexp.MustCompile("^(.*?)($|\\s|[\"'\\x60]|" + regexBoundaries + ")")
	return Token{
		TokenValue: re.FindString(query),
		TokenType:  TOKEN_TYPE_WORD,
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

func Tokenize(query string) (tokens []Token) {
	initialize()

	origLength := len(query)
	oldLength := origLength + 1
	currLength := origLength

	var t Token

	for {
		if currLength < 1 {
			break
		}

		if oldLength <= currLength {
			tokens = append(tokens, Token{
				TokenValue: query,
				TokenType:  TOKEN_TYPE_ERROR,
			})

			return
		}
		oldLength = currLength

		var cacheKey string
		if currLength >= MAX_CACHE_KEY_SIZE {
			cacheKey = query[:MAX_CACHE_KEY_SIZE]
		}

		var tokenLength int
		if cacheKey != "" {
			if cachedToken, exists := tokenCache[cacheKey]; exists {
				t = cachedToken
				tokenLength = len(t.TokenValue)
			}
		}
		if tokenLength < 1 {
			t = getNextToken(query, t)
			tokenLength = len(t.TokenValue)

			if cacheKey != "" && tokenLength < MAX_CACHE_KEY_SIZE {
				if tokenCache == nil {
					tokenCache = make(map[string]Token)
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
