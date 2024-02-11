package genid

import (
	"fmt"

	"github.com/segmentio/ksuid"
)

type Entity string

const (
	EntityRadar         Entity = "rad"
	EntityRadarItem     Entity = "rad_itm"
	EntityRadarQuadrant Entity = "rad_qdt"
)

func New(e Entity) string {
	return fmt.Sprintf("%s_%s", e, ksuid.New().String())
}
