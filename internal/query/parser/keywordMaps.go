package query

var keywords = map[string]TokenType{
	"AND":    LOGICAL,
	"OR":     LOGICAL,
	"RECALL": COMMAND,
}
