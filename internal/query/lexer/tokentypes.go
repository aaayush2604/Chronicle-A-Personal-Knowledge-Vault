package lexer

type TokenType string

const (
	ALL        TokenType = "all"
	COMMAND    TokenType = "command"
	FIELD      TokenType = "field"
	OPERATOR   TokenType = "operator"
	LOGICAL    TokenType = "logical"
	VALUE      TokenType = "value"
	DATE       TokenType = "date"
	TIME       TokenType = "time"
	LBRACKET   TokenType = "left_bracket"
	RBRACKET   TokenType = "right_bracket"
	LPAREN     TokenType = "left_parenthesis"
	RPAREN     TokenType = "right_parenthesis"
	COMMA      TokenType = "comma"
	STRING     TokenType = "string"
	NUMBER     TokenType = "number"
	IDENTIFIER TokenType = "identifier"
	EOF        TokenType = "end_of_input"
)
