package hystrix

import (
	"errors"
)

type (
	command struct {
		circuit  *Circuit
		run      runFunc
		fallback fallbackFunc
		events   []eventType
	}
)

var ErrorCircuitOpen = errors.New("circuit open error")

func (cmd *command) do() {
	defer cmd.reportAllEvents()

	if !cmd.circuit.allowRequest() {
		cmd.report(circuitOpenEvent)

		cmd.tryFallback(ErrorCircuitOpen)
		return
	}

	if err := cmd.run(); err != nil {
		cmd.tryFallback(err)
		return
	}

	cmd.report(successEvent)
}

func (cmd *command) tryFallback(err error) {
	cmd.report(failureEvent)
	if cmd.fallback == nil {
		return
	}

	if fallbackErr := cmd.fallback(err); fallbackErr != nil {
		cmd.report(fallbackFailEvent)
		return
	}
	cmd.report(fallbackSuccessEvent)
}

func (cmd *command) report(event eventType) {
	cmd.events = append(cmd.events, event)
}

func (cmd *command) reportAllEvents() {
	cmd.circuit.reportEvent(cmd.events)
}
