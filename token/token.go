package token

type Token int

const (
	ILLEGAL Token = iota
	EOF

	IDENT
	BYTESLITERAL
	DECLITERAL
	HEXLITERAL
	STRINGLITERAL

	PLUS
	LCURL
	RCURL
	LBRACK
	RBRACK
	LPAREN
	RPAREN
	EQ
	SEMI
	BANG
	LARR
	COMMA

	GEQ
	LEQ
	LT
	GT
	EQEQ
	NEQ
	MINUS
	STAR
	DIV
	DOT
	COLON
	PLUSEQ
	MINUSEQ
	LOGAND
	LOGOR
	BITAND
	BITOR
	BITNEG
	LSHIFT
	RSHIFT
	EXP
	AT

	keyword_beg
	ASSEMBLY
	ASSERT
	CONST
	CONSTRUCTOR
	CONTRACT
	CLASS
	ELSE
	FOR
	FROM
	FUNC
	IF
	IMPORT
	MAP
	POOL
	RETURN
	STORE
	STORAGE
	VAR
	WHILE
	TRUE
	FALSE
	NEW
	CALL
	GAS
	DATA
	VALUE
	MINT
	BURN
	INTERFACE
	PAYABLE
	EVENT
	EMIT
	INDEXED
	AS
	SOLIDITY
	FALLBACK
	keyword_end
)

var Tokens = [...]string{
	ILLEGAL:       "ILLEGAL",
	EOF:           "EOF",
	IDENT:         "IDENT",
	BYTESLITERAL:  "BYTESLITERAL",
	DECLITERAL:    "DECLITERAL",
	HEXLITERAL:    "HEXLITERAL",
	STRINGLITERAL: "STRINGLITERAL",

	PLUS:    "+",
	LCURL:   "{",
	RCURL:   "}",
	LBRACK:  "[",
	RBRACK:  "]",
	LPAREN:  "(",
	RPAREN:  ")",
	EQ:      "=",
	SEMI:    ";",
	BANG:    "!",
	LARR:    "<-",
	COMMA:   ",",
	GEQ:     ">=",
	LEQ:     "<=",
	LT:      "<",
	GT:      ">",
	EQEQ:    "==",
	NEQ:     "!=",
	MINUS:   "-",
	STAR:    "*",
	DIV:     "/",
	DOT:     ".",
	COLON:   ":",
	PLUSEQ:  "+=",
	MINUSEQ: "-=",
	LOGAND:  "&&",
	LOGOR:   "||",
	BITAND:  "&",
	BITOR:   "|",
	BITNEG:  "^",
	LSHIFT:  "<<",
	RSHIFT:  ">>",
	EXP:     "**",
	AT:      "@",

	ASSEMBLY:    "assembly",
	ASSERT:      "assert",
	CONST:       "const",
	CONSTRUCTOR: "constructor",
	CONTRACT:    "contract",
	CLASS:       "class",
	ELSE:        "else",
	FOR:         "for",
	FROM:        "from",
	FUNC:        "func",
	IF:          "if",
	IMPORT:      "import",
	MAP:         "map",
	POOL:        "pool",
	RETURN:      "return",
	STORE:       "store",
	STORAGE:     "storage",
	VAR:         "var",
	WHILE:       "while",
	TRUE:        "true",
	FALSE:       "false",
	NEW:         "new",
	CALL:        "call",
	GAS:         "gas",
	DATA:        "data",
	VALUE:       "value",
	MINT:        "mint",
	BURN:        "burn",
	INTERFACE:   "interface",
	PAYABLE:     "payable",
	EVENT:       "event",
	EMIT:        "emit",
	INDEXED:     "indexed",
	AS:          "as",
	SOLIDITY:    "solidity",
	FALLBACK:    "fallback",
}

var keywords map[string]Token

func init() {
	keywords = make(map[string]Token)
	for i := keyword_beg + 1; i < keyword_end; i++ {
		keywords[Tokens[i]] = i
	}
}

func Lookup(ident string) Token {
	if tok, is_keyword := keywords[ident]; is_keyword {
		return tok
	}
	return IDENT
}
