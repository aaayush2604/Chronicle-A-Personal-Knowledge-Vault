package query

type Field string

const (
	FieldContent   Field = "content"
	FieldTimeStamp Field = "timestamp"
	FieldType      Field = "type"
)

type Operator string

const (
	OpEquals   Operator = "equals"
	OpContains Operator = "contains"
	OpAfter    Operator = "after"
	OpBefore   Operator = "before"
)

type Predicate struct {
	Field    Field
	Operator Operator
	Value    string
}
