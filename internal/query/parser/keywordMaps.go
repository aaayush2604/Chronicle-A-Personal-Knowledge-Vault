package parser

var keywords = map[string]TokenType{
	"AND":    LOGICAL,
	"OR":     LOGICAL,
	"RECALL": COMMAND,
}
