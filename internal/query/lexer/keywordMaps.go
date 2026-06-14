package lexer

var keywords = map[string]TokenType{
	"AND":      LOGICAL,
	"OR":       LOGICAL,
	"RECALL":   COMMAND,
	"REM":      COMMAND,
	"REMEMBER": COMMAND,
	"ALL":      ALL,
}
