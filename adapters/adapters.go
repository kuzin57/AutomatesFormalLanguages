package adapters

import "workspace/internal/automate"

type AutomateAdapter interface {
	Fetch() automate.Automate
}

type automateAdapter struct {
}

func (a *automateAdapter) Fetch() automate.Automate {
	return nil
}
