package automataadapter

import (
	"workspace/adapters"
	"workspace/internal/entities/automata"
)

func NewAutomateAdapter(name string) adapters.AutomataAdapter {
	return &automataAdapter{automata: automata.NewAutomata()}
}
