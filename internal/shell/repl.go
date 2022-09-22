package shell

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"workspace/internal/actions"

	"workspace/internal/shell/shlex"

	"github.com/peterh/liner"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/xerrors"
)

type commandRouter struct {
	shell *Shell

	adapters actions.ActionAdapters

	line *liner.State

	commands *cobra.Command
}

func (r *commandRouter) Close() error {
	if r.line != nil {
		r.writeHistory()
		line := r.line
		r.line = nil
		return line.Close()
	}
	return nil
}

func (r *commandRouter) initLine() {
	r.line.SetCtrlCAborts(true)
	r.line.SetTabCompletionStyle(liner.TabCircular)
	r.line.SetCompleter(r.complete) // Maybe change to SetWordCompleter
	r.readHistory()
}

func (r *commandRouter) readHistory() {
	if hs := r.shell.historyReader(); hs != nil {
		_, err := r.line.ReadHistory(hs)
		r.shell.reportError(err)
		_ = hs.Close()
	}
}

func (r *commandRouter) writeHistory() {
	if hs := r.shell.historyWriter(); hs != nil {
		_, err := r.line.WriteHistory(hs)
		r.shell.reportError(err)
		_ = hs.Close()
	}
}

func (r *commandRouter) command() (err error) {
	command, cmdErr := r.line.Prompt(r.shell.prompt)
	if cmdErr == nil {
		if !strings.HasPrefix(command, " ") { // Skip from history commands started from space.
			r.line.AppendHistory(command)
		}
		return r.processCommand(command)
	} else if cmdErr == liner.ErrPromptAborted {
		return errors.New("Aborted")
	}
	return cmdErr
}

var (
	errQuit           = xerrors.Errorf("quit")
	errUnknownCommand = xerrors.Errorf("unknown command")
)

func (r *commandRouter) processCommand(command string) error {
	args, _, err := r.parseCommand(command)
	if err != nil {
		return fmt.Errorf("command parse error: %w", err)
	}
	if len(args) == 0 {
		return nil
	}
	args = unquoteArgs(args)

	cmd, _, _ := r.commands.Find(args)
	if cmd != nil {
		defer func() {
			cmd.Flags().VisitAll(func(f *pflag.Flag) {
				_ = f.Value.Set(f.DefValue)
				f.Changed = false
			})
		}()
	}
	r.commands.SetArgs(args)

	return r.commands.Execute()
}

func (r *commandRouter) parseCommand(command string) (args []string, poses []int, err error) {
	s := shlex.NewTokenizer(strings.NewReader(command))
	var tok *shlex.Token
	for {
		tok, err = s.Next()

		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}
		tt := tok.Value()
		start := tok.Start()
		args = append(args, tt)
		poses = append(poses, start)
	}
	return args, poses, nil
}

func (r *commandRouter) complete(line string) (result []string) {
	args, poses, err := r.parseCommand(line)
	if err != nil {
		return nil
	}
	switch {
	case len(args) == 0:
		return r.commands.SuggestionsFor(line)
	case strings.HasSuffix(line, " "):
		cmd, cmdArgs, err := r.commands.Find(args)
		if err != nil {
			return nil
		}
		if len(cmdArgs) > 0 {
			// TODO: need flags and params suggest
			return nil
		}
		for _, s := range cmd.SuggestionsFor("") {
			result = append(result, line+s)
		}
		return
	default:
		cmd, cmdArgs, err := r.commands.Find(args[:len(args)-1])
		if err != nil {
			return nil
		}
		if len(cmdArgs) > 0 {
			// TODO: need flags and params suggest
			return nil
		}
		p := poses[len(args)-1]
		prefix := string([]byte(line)[:p])
		for _, s := range cmd.SuggestionsFor(args[len(args)-1]) {
			result = append(result, prefix+s)
		}
		return
	}
}

func unquoteArgs(in []string) (out []string) {
	out = make([]string, len(in))
	for i, v := range in {
		out[i] = v
		if val, err := strconv.Unquote(v); err == nil {
			out[i] = val
		}
	}
	return
}
