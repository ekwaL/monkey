package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

func (t Token) String() string {
	return "Token{Type: " + string(t.Type) + ", Literal: " + t.Literal + "}"
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"
	STRING     = "STRING"

	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	STAR   = "*"
	SLASH  = "/"
	BANG   = "!"

	OR  = "||"
	AND = "&&"

	LESS          = "<"
	GREATER       = ">"
	EQUAL_EQUAL   = "=="
	NOT_EQUAL     = "!="
	LESS_EQUAL    = "<="
	GREATER_EQUAL = ">="

	DOT       = "."
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"
	LHASHBRACE = "{|"
	RHASHBRACE = "|}"

	// Keywords
	CLASS    = "CLASS"
	THIS     = "THIS"
	SUPER    = "SUPER"
	FUNCTION = "FUNCTION"
	LET      = "LET"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"
)

const (
	THIS_KEYWORD        = "this"
	SUPER_KEYWORD       = "super"
	INITIALIZER_KEYWORD = "init"
)

var keywords = map[string]TokenType{
	"let":         LET,
	"class":       CLASS,
	THIS_KEYWORD:  THIS,
	SUPER_KEYWORD: SUPER,
	"fn":          FUNCTION,
	"return":      RETURN,
	"if":          IF,
	"else":        ELSE,
	"true":        TRUE,
	"false":       FALSE,
	"null":        NULL,
}

func LookupKeyword(identifier string) TokenType {
	if token, ok := keywords[identifier]; ok {
		return token
	}

	return IDENTIFIER
}
