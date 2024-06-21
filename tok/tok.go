package tok

type TokenType string

type Tok struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT     = "IDENT"
	INT       = "INT"
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	ASTERISK  = "*"
	SLASH     = "/"
	SEMICOLON = ";"
	GT        = ">"
	LT        = "<"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	VAR   = "VAR"
	IF    = "IF"
	ELSE  = "ELSE"
	TRUE  = "TRUE"
	FALSE = "FALSE"
	EXIT  = "EXIT"
)

var keywords = map[string]TokenType{
	"var":   VAR,
	"if":    IF,
	"else":  ELSE,
	"true":  TRUE,
	"false": FALSE,
	"exit":  EXIT,
}

func KeywordLookup(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
