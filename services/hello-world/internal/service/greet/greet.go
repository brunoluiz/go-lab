package greet

import (
	"errors"

	"golang.org/x/text/language"
)

var ErrNotImplemented = errors.New("not implemented")

type Greeter struct{}

func New() *Greeter {
	return &Greeter{}
}

func (g *Greeter) Hello(l string) (string, error) {
	lang, err := language.Parse(l)
	if err != nil {
		return "", err
	}

	switch lang {
	case language.English:
		return "hello", nil
	case language.Portuguese:
		return "olá", nil
	case language.German:
		return "hallo", nil
	default:
		return "", ErrNotImplemented
	}
}
