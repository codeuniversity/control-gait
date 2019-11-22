package gait

// Action ...
type Action struct {
	started bool

	legs    []*Leg
	command string

	next *Action
}

// Leg ...
type Leg struct {
	name string
	done bool
}

var movement = []string{"move up", "rotate 30", "move down", "rotate -30"}

// NewMovement returns an action chain
func NewMovement(legsNames []string) Movement {
	var firstAction *Action
	var previousAction *Action

	for _, command := range movement {

		legs := []*Leg{}
		for _, name := range legsNames {
			legs = append(legs, &Leg{name: name})
		}

		action := &Action{
			legs:    legs,
			command: command,
		}
		if previousAction == nil {
			firstAction = action
		} else {
			previousAction.next = action
		}
		previousAction = action
	}

	return firstAction
}

// Done ...
func (a *Action) Done() bool {
	for _, leg := range a.legs {
		if !leg.done {
			return false
		}
	}

	return true
}

// Movement is an action chain
type Movement *Action
