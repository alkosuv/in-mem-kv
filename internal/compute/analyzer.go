package compute

import (
	"regexp"
	"strings"
)

const (
	setTokensNumber int = 3
	getTokensNumber int = 2
	delTokensNumber int = 2
)

var regexpKey = regexp.MustCompile(`^[a-zA-Z0-9/*_.]+$`)

// analyzeQuery - метод проверяет можно ли обработать запрос
func analyzeQuery(tokens tokens) error {
	if len(tokens) == 0 {
		return ErrInvalidQuery
	}

	switch strings.ToUpper(tokens[commandTokenIndex]) {
	case SetCommand.String():
		if len(tokens) != setTokensNumber {
			return ErrInvalidNumberArgument
		}
	case GetCommand.String():
		if len(tokens) != getTokensNumber {
			return ErrInvalidNumberArgument
		}
	case DelCommand.String():
		if len(tokens) != delTokensNumber {
			return ErrInvalidNumberArgument
		}
	default:
		return ErrInvalidCommand
	}

	if !regexpKey.MatchString(tokens[keyTokenIndex]) {
		return ErrInvalidSymbol
	}

	return nil
}
