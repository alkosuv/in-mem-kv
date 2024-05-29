package compute

import "errors"

var ErrInvalidQuery = errors.New("invalid query")
var ErrInvalidSymbol = errors.New("the query contains an invalid symbol")
var ErrInvalidCommand = errors.New("the query has an invalid command")
var ErrInvalidNumberArgument = errors.New("the query has an invalid number arguments")
