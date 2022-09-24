package automatesadapters

import (
	"workspace/adapters"
	"workspace/internal/automate"
	"workspace/internal/config"
)

func NewAutomateAdapter(cfg config.AdaptersConfig) adapters.AutomateAdapter {
	if cfg.IsDeterministic {
		return &faAutomateAdapter{automate: automate.NewFA()}
	}
	return &nfaAutomateAdapter{automate: automate.NewNFA()}
}
