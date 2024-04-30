package genid

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

func New[T comparable](e T) string {
	return fmt.Sprintf("%s_%s", e, ksuid.New().String())
}
