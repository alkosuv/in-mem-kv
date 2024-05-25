package compute

type tokens []string

type tokenIndex int

const (
	commandTokenIndex tokenIndex = iota
	keyTokenIndex
	valueTokenIndex
)
