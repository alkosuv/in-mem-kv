package compute

type Command string

func (c Command) String() string {
	return string(c)
}

const (
	SetCommand Command = "SET"
	GetCommand Command = "GET"
	DelCommand Command = "DEL"
)

// StringToCommand преобразование строки в команду
var StringToCommand = map[string]Command{
	"SET": SetCommand,
	"GET": GetCommand,
	"DEL": DelCommand,
}
