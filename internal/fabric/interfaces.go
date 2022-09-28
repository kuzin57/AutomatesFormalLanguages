package fabric

import "workspace/internal/entities/automate"

type AutomateFabric interface {
	Create() (*automate.Automate, error)
}
