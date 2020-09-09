package cli

import (
    "errors"
)

var (
    ErrCommandNoMatch = errors.New("Arguments do not match the command")
)


type Command interface {
    ShortDescription() string
    Execute() error
    Parse(args []string) error
}





