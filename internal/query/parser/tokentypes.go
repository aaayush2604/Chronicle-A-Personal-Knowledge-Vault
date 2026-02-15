package query

type TokenType string

const (
	COMMAND  TokenType = "command"
	FIELD    TokenType = "field"
	OPERATOR TokenType = "operator"
	LOGICAL  TokenType = "logical"
	VALUE    TokenType = "value"
	DATE     TokenType = "date"
	TIME     TokenType = "time"
	LBRACKET TokenType = "left_bracket"
	RBRACKET TokenType = "right_bracket"
	COMMA    TokenType = "comma"
	EOF      TokenType = "end_of_input"
)
