package workspaces

import (
	"errors"
	"fmt"
	"os"

	"github.com/reecerussell/passport"
)

// RunScriptCommand is a command used to execute a script.
var RunScriptCommand = &passport.Command{
	Name:        "run",
	Description: "used execute a script in a workspace",
	Execute: func(cmd *passport.Command, ctx *passport.CommandContext) error {
		cnf, err := passport.LoadConfig(ctx.ConfigDir, ctx.Fs)
		if err != nil {
			return err
		}

		wd, _ := os.Getwd()
		w, err := cnf.GetWorkspace(wd)
		if err != nil {
			return err
		}

		if len(os.Args) < 3 {
			return errors.New("run: no script name specified")
		}

		name := os.Args[2]
		s, err := w.GetScript(name)
		if err != nil {
			return err
		}

		exitCode, err := s.Run(ctx.Crypto)
		if err != nil {
			return err
		}

		fmt.Printf("Exited with code %d\n", exitCode)

		return nil
	},
}
