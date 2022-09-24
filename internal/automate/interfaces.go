package automate

type Automate interface {
	DeleteEps() error
	Read(string) bool
	Check() bool
	AddNewWord(string) error
	Cycle() error
	Concat(Automate) error
	Join(Automate) error
}
