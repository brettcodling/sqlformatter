package tokens

import (
	"regexp"
	"strings"
)

var (
	boundaries = []string{
		",", ";", ":", ")", "(", ".", "=", "<", ">", "+", "-", "*", "/", "!", "^", "%", "|", "&", "#",
	}
	functions = []string{
		"ABS", "ACOS", "ADDDATE", "ADDTIME", "AES_DECRYPT", "AES_ENCRYPT", "AREA", "ASBINARY", "ASCII", "ASIN",
		"ASTEXT", "ATAN", "ATAN2", "AVG", "BDMPOLYFROMTEXT", "BDMPOLYFROMWKB", "BDPOLYFROMTEXT", "BDPOLYFROMWKB",
		"BENCHMARK", "BIN", "BIT_AND", "BIT_COUNT", "BIT_LENGTH", "BIT_OR", "BIT_XOR", "BOUNDARY", "BUFFER", "CAST",
		"CEIL", "CEILING", "CENTROID", "CHAR", "CHARACTER_LENGTH", "CHARSET", "CHAR_LENGTH", "COALESCE",
		"COERCIBILITY", "COLLATION", "COMPRESS", "CONCAT", "CONCAT_WS", "CONNECTION_ID", "CONTAINS", "CONV", "CONVERT",
		"CONVERT_TZ", "CONVEXHULL", "COS", "COT", "COUNT", "CRC32", "CROSSES", "CURDATE", "CURRENT_DATE",
		"CURRENT_TIME", "CURRENT_TIMESTAMP", "CURRENT_USER", "CURTIME", "DATABASE", "DATE", "DATEDIFF", "DATE_ADD",
		"DATE_DIFF", "DATE_FORMAT", "DATE_SUB", "DAY", "DAYNAME", "DAYOFMONTH", "DAYOFWEEK", "DAYOFYEAR", "DECODE",
		"DEFAULT", "DEGREES", "DES_DECRYPT", "DES_ENCRYPT", "DIFFERENCE", "DIMENSION", "DISJOINT", "DISTANCE", "ELT",
		"ENCODE", "ENCRYPT", "ENDPOINT", "ENVELOPE", "EQUALS", "EXP", "EXPORT_SET", "EXTERIORRING", "EXTRACT",
		"EXTRACTVALUE", "FIELD", "FIND_IN_SET", "FLOOR", "FORMAT", "FOUND_ROWS", "FROM_DAYS", "FROM_UNIXTIME",
		"GEOMCOLLFROMTEXT", "GEOMCOLLFROMWKB", "GEOMETRYCOLLECTION", "GEOMETRYCOLLECTIONFROMTEXT",
		"GEOMETRYCOLLECTIONFROMWKB", "GEOMETRYFROMTEXT", "GEOMETRYFROMWKB", "GEOMETRYN", "GEOMETRYTYPE",
		"GEOMFROMTEXT", "GEOMFROMWKB", "GET_FORMAT", "GET_LOCK", "GLENGTH", "GREATEST", "GROUP_CONCAT",
		"GROUP_UNIQUE_USERS", "HEX", "HOUR", "IF", "IFNULL", "INET_ATON", "INET_NTOA", "INSERT", "INSTR",
		"INTERIORRINGN", "INTERSECTION", "INTERSECTS", "INTERVAL", "ISCLOSED", "ISEMPTY", "ISNULL", "ISRING",
		"ISSIMPLE", "IS_FREE_LOCK", "IS_USED_LOCK", "LAST_DAY", "LAST_INSERT_ID", "LCASE", "LEAST", "LEFT", "LENGTH",
		"LINEFROMTEXT", "LINEFROMWKB", "LINESTRING", "LINESTRINGFROMTEXT", "LINESTRINGFROMWKB", "LN", "LOAD_FILE",
		"LOCALTIME", "LOCALTIMESTAMP", "LOCATE", "LOG", "LOG10", "LOG2", "LOWER", "LPAD", "LTRIM", "MAKEDATE",
		"MAKETIME", "MAKE_SET", "MASTER_POS_WAIT", "MAX", "MBRCONTAINS", "MBRDISJOINT", "MBREQUAL", "MBRINTERSECTS",
		"MBROVERLAPS", "MBRTOUCHES", "MBRWITHIN", "MD5", "MICROSECOND", "MID", "MIN", "MINUTE", "MLINEFROMTEXT",
		"MLINEFROMWKB", "MOD", "MONTH", "MONTHNAME", "MPOINTFROMTEXT", "MPOINTFROMWKB", "MPOLYFROMTEXT",
		"MPOLYFROMWKB", "MULTILINESTRING", "MULTILINESTRINGFROMTEXT", "MULTILINESTRINGFROMWKB", "MULTIPOINT",
		"MULTIPOINTFROMTEXT", "MULTIPOINTFROMWKB", "MULTIPOLYGON", "MULTIPOLYGONFROMTEXT", "MULTIPOLYGONFROMWKB",
		"NAME_CONST", "NULLIF", "NUMGEOMETRIES", "NUMINTERIORRINGS", "NUMPOINTS", "OCT", "OCTET_LENGTH",
		"OLD_PASSWORD", "ORD", "OVERLAPS", "PASSWORD", "PERIOD_ADD", "PERIOD_DIFF", "PI", "POINT", "POINTFROMTEXT",
		"POINTFROMWKB", "POINTN", "POINTONSURFACE", "POLYFROMTEXT", "POLYFROMWKB", "POLYGON", "POLYGONFROMTEXT",
		"POLYGONFROMWKB", "POSITION", "POW", "POWER", "QUARTER", "QUOTE", "RADIANS", "RAND", "RELATED", "RELEASE_LOCK",
		"REPEAT", "REPLACE", "REVERSE", "RIGHT", "ROUND", "ROW_COUNT", "RPAD", "RTRIM", "SCHEMA", "SECOND",
		"SEC_TO_TIME", "SESSION_USER", "SHA", "SHA1", "SIGN", "SIN", "SLEEP", "SOUNDEX", "SPACE", "SQRT", "SRID",
		"STARTPOINT", "STD", "STDDEV", "STDDEV_POP", "STDDEV_SAMP", "STRCMP", "STR_TO_DATE", "SUBDATE", "SUBSTR",
		"SUBSTRING", "SUBSTRING_INDEX", "SUBTIME", "SUM", "SYMDIFFERENCE", "SYSDATE", "SYSTEM_USER", "TAN", "TIME",
		"TIMEDIFF", "TIMESTAMP", "TIMESTAMPADD", "TIMESTAMPDIFF", "TIME_FORMAT", "TIME_TO_SEC", "TOUCHES", "TO_DAYS",
		"TRIM", "TRUNCATE", "UCASE", "UNCOMPRESS", "UNCOMPRESSED_LENGTH", "UNHEX", "UNIQUE_USERS", "UNIX_TIMESTAMP",
		"UPDATEXML", "UPPER", "USER", "UTC_DATE", "UTC_TIME", "UTC_TIMESTAMP", "UUID", "VARIANCE", "VAR_POP",
		"VAR_SAMP", "VERSION", "WEEK", "WEEKDAY", "WEEKOFYEAR", "WITHIN", "X", "Y", "YEAR", "YEARWEEK",
	}
	maxCacheKeySize       = 15
	regexBoundaries       string
	regexFunction         string
	regexReserved         string
	regexReservedToplevel string
	regexReservedNewline  string
	reserved              = []string{
		"GEOMETRYCOLLECTIONFROMTEXT", "GEOMETRYCOLLECTIONFROMWKB", "MULTILINESTRINGFROMTEXT", "MULTILINESTRINGFROMWKB",
		"MULTIPOLYGONFROMTEXT", "MULTIPOLYGONFROMWKB", "UNCOMPRESSED_LENGTH", "GEOMETRYCOLLECTION",
		"GROUP_UNIQUE_USERS", "LINESTRINGFROMTEXT", "MULTIPOINTFROMTEXT", "CURRENT_TIMESTAMP", "LINESTRINGFROMWKB",
		"MULTIPOINTFROMWKB", "CHARACTER_LENGTH", "GEOMCOLLFROMTEXT", "GEOMETRYFROMTEXT", "NUMINTERIORRINGS",
		"BDMPOLYFROMTEXT", "GEOMCOLLFROMWKB", "GEOMETRYFROMWKB", "MASTER_POS_WAIT", "MULTILINESTRING",
		"POLYGONFROMTEXT", "SUBSTRING_INDEX", "BDMPOLYFROMWKB", "BDPOLYFROMTEXT", "LAST_INSERT_ID", "LOCALTIMESTAMP",
		"MPOINTFROMTEXT", "POINTONSURFACE", "POLYGONFROMWKB", "UNIX_TIMESTAMP", "BDPOLYFROMWKB", "CONNECTION_ID",
		"FROM_UNIXTIME", "INTERIORRINGN", "MBRINTERSECTS", "MLINEFROMTEXT", "MPOINTFROMWKB", "MPOLYFROMTEXT",
		"NUMGEOMETRIES", "POINTFROMTEXT", "SYMDIFFERENCE", "TIMESTAMPDIFF", "UTC_TIMESTAMP", "COERCIBILITY",
		"CURRENT_DATE", "CURRENT_TIME", "CURRENT_USER", "EXTERIORRING", "EXTRACTVALUE", "GEOMETRYTYPE", "GEOMFROMTEXT",
		"GROUP_CONCAT", "INTERSECTION", "IS_FREE_LOCK", "IS_USED_LOCK", "LINEFROMTEXT", "MLINEFROMWKB", "MPOLYFROMWKB",
		"MULTIPOLYGON", "OCTET_LENGTH", "OLD_PASSWORD", "POINTFROMWKB", "POLYFROMTEXT", "RELEASE_LOCK", "SESSION_USER",
		"TIMESTAMPADD", "UNIQUE_USERS", "AES_DECRYPT", "AES_ENCRYPT", "CHAR_LENGTH", "DATE_FORMAT", "DES_DECRYPT",
		"DES_ENCRYPT", "FIND_IN_SET", "GEOMFROMWKB", "LINEFROMWKB", "MBRCONTAINS", "MBRDISJOINT", "MBROVERLAPS",
		"MICROSECOND", "PERIOD_DIFF", "POLYFROMWKB", "SEC_TO_TIME", "STDDEV_SAMP", "STR_TO_DATE", "SYSTEM_USER",
		"TIME_FORMAT", "TIME_TO_SEC", "BIT_LENGTH", "CONVERT_TZ", "CONVEXHULL", "DAYOFMONTH", "DIFFERENCE",
		"EXPORT_SET", "FOUND_ROWS", "GET_FORMAT", "INTERSECTS", "LINESTRING", "MBRTOUCHES", "MULTIPOINT", "NAME_CONST",
		"PERIOD_ADD", "STARTPOINT", "STDDEV_POP", "UNCOMPRESS", "WEEKOFYEAR", "BENCHMARK", "BIT_COUNT", "COLLATION",
		"CONCAT_WS", "DATE_DIFF", "DAYOFWEEK", "DAYOFYEAR", "DIMENSION", "FROM_DAYS", "GEOMETRYN", "INET_ATON",
		"INET_NTOA", "LOAD_FILE", "LOCALTIME", "MBRWITHIN", "MONTHNAME", "NUMPOINTS", "ROW_COUNT", "SUBSTRING",
		"TIMESTAMP", "UPDATEXML", "ASBINARY", "BOUNDARY", "CENTROID", "COALESCE", "COMPRESS", "CONTAINS", "DATABASE",
		"DATEDIFF", "DATE_ADD", "DATE_SUB", "DISJOINT", "DISTANCE", "ENDPOINT", "ENVELOPE", "GET_LOCK", "GREATEST",
		"INTERVAL", "ISCLOSED", "ISSIMPLE", "LAST_DAY", "MAKEDATE", "MAKETIME", "MAKE_SET", "MBREQUAL", "OVERLAPS",
		"PASSWORD", "POSITION", "TIMEDIFF", "TRUNCATE", "UTC_DATE", "UTC_TIME", "VARIANCE", "VAR_SAMP", "YEARWEEK",
		"ADDDATE", "ADDTIME", "BIT_AND", "BIT_XOR", "CEILING", "CHARSET", "CONVERT", "CROSSES", "CURDATE", "CURTIME",
		"DAYNAME", "DEFAULT", "DEGREES", "ENCRYPT", "EXTRACT", "GLENGTH", "ISEMPTY", "POLYGON", "QUARTER", "RADIANS",
		"RELATED", "REPLACE", "REVERSE", "SOUNDEX", "SUBDATE", "SUBTIME", "SYSDATE", "TOUCHES", "TO_DAYS", "VAR_POP",
		"VERSION", "WEEKDAY", "ASTEXT", "BIT_OR", "BUFFER", "CONCAT", "DECODE", "ENCODE", "EQUALS", "FORMAT", "IFNULL",
		"INSERT", "ISNULL", "ISRING", "LENGTH", "LOCATE", "MINUTE", "NULLIF", "POINTN", "REPEAT", "SCHEMA", "SECOND",
		"STDDEV", "STRCMP", "SUBSTR", "WITHIN", "ASCII", "ATAN2", "COUNT", "CRC32", "FIELD", "FLOOR", "INSTR", "LCASE",
		"LEAST", "LOG10", "LOWER", "LTRIM", "MONTH", "POINT", "POWER", "QUOTE", "RIGHT", "ROUND", "RTRIM", "SLEEP",
		"SPACE", "UCASE", "UNHEX", "UPPER", "ACOS", "AREA", "ASIN", "ATAN", "CAST", "CEIL", "CHAR", "CONV", "DATE",
		"HOUR", "LEFT", "LOG2", "LPAD", "RAND", "RPAD", "SHA1", "SIGN", "SQRT", "SRID", "TIME", "TRIM", "USER", "UUID",
		"WEEK", "YEAR", "ABS", "AVG", "BIN", "COS", "COT", "DAY", "ELT", "EXP", "HEX", "LOG", "MAX", "MD5", "MID",
		"MIN", "MOD", "OCT", "ORD", "POW", "SHA", "SIN", "STD", "SUM", "TAN", "IF", "LN", "PI", "X", "Y",
	}
	reservedNewline = []string{
		"LEFT OUTER JOIN", "RIGHT OUTER JOIN", "LEFT JOIN", "RIGHT JOIN", "OUTER JOIN", "INNER JOIN", "JOIN", "XOR",
		"OR", "AND",
	}
	reservedToplevel = []string{
		"SELECT", "FROM", "WHERE", "SET", "ORDER BY", "GROUP BY", "LIMIT", "DROP", "VALUES", "UPDATE", "HAVING", "ADD",
		"AFTER", "ALTER TABLE", "DELETE FROM", "UNION ALL", "UNION", "EXCEPT", "INTERSECT",
	}
	tokenCache map[Query]Token
)

type Token struct {
	TokenValue string
	TokenType  int
	Index      int
}

func getNextToken(query Query, prevToken Token) Token {
	// Whitespace
	if token := query.matchWhitespace(); token.TokenValue != "" {
		return token
	}
	// Comment
	if token := query.matchComment(); token.TokenValue != "" {
		return token
	}
	// Quoted
	if token := query.matchQuote(); token.TokenValue != "" {
		return token
	}
	// User defined
	if token := query.matchUserDefined(); token.TokenValue != "" {
		return token
	}
	// Number
	if token := query.matchNumber(); token.TokenValue != "" {
		return token
	}
	// Boundary
	if token := query.matchBoundary(); token.TokenValue != "" {
		return token
	}
	// Reserved from prefixed by '.'
	if token := query.matchReserved(prevToken); token.TokenValue != "" {
		return token
	}
	// Function suffixed by '('
	if token := query.matchFunction(); token.TokenValue != "" {
		return token
	}
	// Non reserved
	return query.matchNonReserved()
}

func joinQuotedRegexp(values []string) string {
	var quoted []string
	for _, value := range values {
		quoted = append(quoted, regexp.QuoteMeta(value))
	}

	return "(" + strings.Join(quoted, "|") + ")"
}

func Tokenize(query Query) (tokens []Token) {
	if regexBoundaries == "" {
		regexBoundaries = joinQuotedRegexp(boundaries)
		regexReserved = joinQuotedRegexp(reserved)
		regexReservedToplevel = strings.ReplaceAll(joinQuotedRegexp(reservedToplevel), " ", "\\s+")
		regexReservedNewline = strings.ReplaceAll(joinQuotedRegexp(reservedNewline), " ", "\\s+")
		regexFunction = joinQuotedRegexp(functions)
	}

	oldLength := len(query) + 1
	currLength := len(query)

	var t Token

	for {
		if currLength < 1 {
			break
		}

		if oldLength <= currLength {
			tokens = append(tokens, Token{
				TokenValue: string(query),
				TokenType:  Token_type_error,
			})

			return
		}
		oldLength = currLength

		var cacheKey Query
		if currLength >= maxCacheKeySize {
			cacheKey = query[:maxCacheKeySize]
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

			if cacheKey != "" && tokenLength < maxCacheKeySize {
				if tokenCache == nil {
					tokenCache = make(map[Query]Token)
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
