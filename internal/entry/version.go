package entry

type SchemaVersion int

const (
	Version1 SchemaVersion = 1 //bare golang based entry and retrieval

	Version2 SchemaVersion = 2 //added query language for recalling thoughts

	Version3 SchemaVersion = 3 //added query language support for adding thoughts with tags

	CurrentVersion = Version3
)
