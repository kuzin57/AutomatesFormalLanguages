package shell

import (
	"fmt"
	"io"
	"os"
	"workspace/internal/automate"

	"github.com/mitchellh/go-homedir"
	"github.com/peterh/liner"
	"github.com/spf13/cobra"
	"golang.org/x/xerrors"
)

type Shell struct {
	router   *commandRouter
	Automate automate.Automate
	prompt   string
}

func (s *Shell) Run() (err error) {
	s.router.line = liner.NewLiner()
	defer func() {
		closeErr := s.router.Close()
		if err == nil {
			err = closeErr
		}
	}()
	s.router.initLine()

	for {
		err := s.router.command()
		switch {
		case xerrors.Is(err, errQuit) || xerrors.Is(err, io.EOF):
			s.reportError(err)
			return nil
		case xerrors.Is(err, errUnknownCommand):
			s.reportError(err)
			s.print("use help for list of available commands")
		case err != nil:
			s.reportError(err)
		}
	}
}

func (s *Shell) Init() error {
	s.prompt = fmt.Sprintf(promptTemplate)

	s.router = &commandRouter{
		shell: s,
	}

	s.router.commands = &cobra.Command{}
	s.router.commands.AddCommand(createCmd)

	registerCreateSubcommands(s)

	return nil
}

func (s *Shell) historyWriter() io.WriteCloser {
	path, err := homedir.Expand(historyFile)
	if err != nil {
		s.reportError(err)
		return nil
	}
	f, err := os.Create(path)
	if err != nil {
		s.reportError(err)
		return nil
	}
	return f
}

func (s *Shell) historyReader() io.ReadCloser {
	path, err := homedir.Expand(historyFile)
	if err != nil {
		s.reportError(err)
		return nil
	}
	f, err := os.Open(path)
	if err != nil {
		s.reportError(err)
		return nil
	}
	return f
}

func (s *Shell) reportError(err error) {
	if err == nil {
		return
	}
	s.print(err.Error())
}

func (s *Shell) print(value string) {
	fmt.Println("  ", value)
}
