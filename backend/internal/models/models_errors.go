package models

import (
	"fmt"
)

type InvalidCommand struct {
	Requested string
}

func (i InvalidCommand) Error() string {
	return fmt.Sprintf("Invalid command: %v", i.Requested)
}

type InvalidStateChange struct {
	CurrentState string
	Requested    string
}

func (i InvalidStateChange) Error() string {
	return fmt.Sprintf("Invalid state change from %v to %v", i.CurrentState, i.Requested)
}

type InvalidArduinoCommand struct {
	CurrentState string
	Requested    string
}

func (i InvalidArduinoCommand) Error() string {
	return fmt.Sprintf("Invalid Arduino Command: %q, current state %v", i.Requested, i.CurrentState)
}
