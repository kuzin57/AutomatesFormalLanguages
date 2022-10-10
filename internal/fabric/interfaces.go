package fabric

import automata "workspace/internal/entities/automata"

type AutomataFabric interface {
	Create() (*automata.Automata, error)
}
