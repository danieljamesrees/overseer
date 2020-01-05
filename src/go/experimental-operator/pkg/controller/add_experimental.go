package controller

import (
	"experimental-operator/pkg/controller/experimental"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, experimental.Add)
}
