package lexer

var keywords = map[string]TokenType{
	"AND":       LOGICAL,
	"OR":        LOGICAL,
	"RECALL":    COMMAND,
	"REM":       COMMAND,
	"REMEMBER":  COMMAND,
	"N":         ETYPE,
	"NOTE":      ETYPE,
	"L":         ETYPE,
	"LEARNING":  ETYPE,
	"Q":         ETYPE,
	"QUESTION":  ETYPE,
	"I":         ETYPE,
	"IDEA":      ETYPE,
	"IMP":       ETYPE,
	"IMPORTANT": ETYPE,
	"ALL":       ALL,
}
