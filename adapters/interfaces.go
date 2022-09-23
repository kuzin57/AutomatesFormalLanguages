package adapters

import (
	"workspace/internal/automate"
)

type AutomateAdapter interface {
	Get() (automate.Automate, error)
	Create(string, []string) error
	AddStar() error
	Join(AutomateAdapter) error
}
