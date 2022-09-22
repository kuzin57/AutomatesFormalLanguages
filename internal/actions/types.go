package actions

import "workspace/adapters"

type ActionAdapters interface {
	Automate() adapters.AutomateAdapter
}
