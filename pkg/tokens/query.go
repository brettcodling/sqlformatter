package tokens

import (
	"regexp"
	"strings"
)

const (
	Token_type_whitespace = iota
	Token_type_word
	Token_type_quote
	Token_type_backtick_quote
	Token_type_reserved
	Token_type_reserved_toplevel
	Token_type_reserved_newline
	Token_type_boundary
	Token_type_comment
	Token_type_block_comment
	Token_type_number
	Token_type_error
	Token_type_variable
)

type Query string

func getQuotedString(query string) (quoted string) {
	re := regexp.MustCompile("^(((\\x60[^\\x60]*($|\\x60))+)|((\\[[^\\]]*($|\\]))(\\][^\\]]*($|\\]))*)|((\"[^\"\\\\]*(?:\\\\.[^\"\\\\]*)*(\"|$))+)|(('[^'\\\\]*(?:\\\\.[^'\\\\]*)*('|$))+))")
	quoted = re.FindString(query)

	return
}

func (q Query) matchBoundary() Token {
	re := regexp.MustCompile("^" + regexBoundaries)
	return Token{
		TokenValue: re.FindString(string(q)),
		TokenType:  Token_type_boundary,
	}
}

func (q Query) matchComment() Token {
	if q[:1] == "#" || (len(q) > 1 && (q[:2] == "--" || q[:2] == "/*")) {
		var last, tokenType int
		if q[:1] == "#" || q[:1] == "-" {
			last = strings.Index(string(q), "\n")
			tokenType = Token_type_comment
		} else {
			last = strings.Index(string(q[2:]), "*/") + 2
			tokenType = Token_type_block_comment
		}

		if last == -1 {
			last = len(q)
		}

		return Token{
			TokenValue: string(q[:last]),
			TokenType:  tokenType,
		}
	}

	return Token{}
}

func (q Query) matchFunction() Token {
	upper := strings.ToUpper(string(q))
	re := regexp.MustCompile("^(" + regexReserved + "[(]|\\s|[)])")
	if match := re.FindString(upper); match != "" {
		return Token{
			TokenValue: string(q[:len(match)-1]),
			TokenType:  Token_type_reserved,
		}
	}

	return Token{}
}

func (q Query) matchNonReserved() Token {
	re := regexp.MustCompile("^(.*?)($|\\s|[\"'\\x60]|" + regexBoundaries + ")")
	return Token{
		TokenValue: re.FindString(string(q)),
		TokenType:  Token_type_word,
	}
}

func (q Query) matchNumber() Token {
	re := regexp.MustCompile("^([0-9]+(\\.[0-9]+)?|0x[0-9a-fA-F]+|0b[01]+)($|\\s|\"'\\x60|" + regexBoundaries + ")")
	return Token{
		TokenValue: re.FindString(string(q)),
		TokenType:  Token_type_number,
	}
}

func (q Query) matchQuote() Token {
	if q[:1] == "\"" || q[:1] == "'" || q[:1] == "`" || q[:1] == "[" {
		tokenType := Token_type_quote
		if q[:1] == "`" || q[:1] == "[" {
			tokenType = Token_type_backtick_quote
		}

		return Token{
			TokenValue: getQuotedString(string(q)),
			TokenType:  tokenType,
		}
	}

	return Token{}
}

func (q Query) matchReserved(prevToken Token) Token {
	upper := strings.ToUpper(string(q))
	if prevToken.TokenValue == "" || prevToken.TokenValue != "." {
		re := regexp.MustCompile("^" + regexReservedToplevel + "($|\\s|" + regexBoundaries + ")")
		if match := re.FindString(upper); match != "" {
			return Token{
				TokenValue: string(q[:len(match)]),
				TokenType:  Token_type_reserved_toplevel,
			}
		}
		re = regexp.MustCompile("^" + regexReservedNewline + "($|\\s|" + regexBoundaries + ")")
		if match := re.FindString(upper); match != "" {
			return Token{
				TokenValue: string(q[:len(match)]),
				TokenType:  Token_type_reserved_newline,
			}
		}
		re = regexp.MustCompile("^" + regexReserved + "($|\\s|" + regexBoundaries + ")")
		if match := re.FindString(upper); match != "" {
			return Token{
				TokenValue: string(q[:len(match)]),
				TokenType:  Token_type_reserved,
			}
		}
	}

	return Token{}
}

func (q Query) matchUserDefined() Token {
	if len(q) > 1 && (q[:1] == "@" || q[:1] == ":") {
		if q[1:2] == "\"" || q[1:2] == "'" || q[1:2] == "`" {
			return Token{
				TokenValue: string(q[:1]) + getQuotedString(string(q[1:])),
				TokenType:  Token_type_variable,
			}
		} else {
			re := regexp.MustCompile("^(" + string(q[:1]) + "[a-zA-Z0-9\\._\\$]+)")
			return Token{
				TokenValue: re.FindString(string(q)),
				TokenType:  Token_type_variable,
			}
		}
	}

	return Token{}
}

func (q Query) matchWhitespace() Token {
	re := regexp.MustCompile("^\\s+")
	return Token{
		TokenValue: re.FindString(string(q)),
		TokenType:  Token_type_whitespace,
	}
}
