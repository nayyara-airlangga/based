package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keywords map[string]TokenType = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LookupType(ident string) TokenType {
	if tokType, ok := keywords[ident]; ok {
		return tokType
	}
	return IDENT
}
