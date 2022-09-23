package automatesadapters

import (
	"workspace/adapters"
	"workspace/internal/config"
)

func NewAutomateAdapter(cfg config.AdaptersConfig) adapters.AutomateAdapter {
	if cfg.IsDeterministic {
		return &faAutomateAdapter{}
	}
	return &nfaAutomateAdapter{}
}
