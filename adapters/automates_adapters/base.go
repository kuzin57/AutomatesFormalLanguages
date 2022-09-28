package automateadapter

import (
	"workspace/adapters"
	"workspace/internal/entities/automate"
)

func NewAutomateAdapter(name string) adapters.AutomateAdapter {
	return &automateAdapter{automate: automate.NewAutomate()}
}
