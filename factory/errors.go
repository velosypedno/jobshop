package factory

import (
	"fmt"
	"strings"

	"github.com/velosypedno/resource-allocation/base"
)

type MissingMachinesError struct {
	MissingTypes []base.MachineType
}

func (e *MissingMachinesError) Error() string {
	types := make([]string, len(e.MissingTypes))
	for _, t := range e.MissingTypes {
		types = append(types, t.String())
	}
	return fmt.Sprintf("missing required machine types: %s", strings.Join(types, ", "))
}
