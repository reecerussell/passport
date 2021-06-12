package workspaces

import (
	"fmt"
	"os"

	"github.com/reecerussell/passport"
)

// Command is the main entrypoint command for operations around workspace scripts.
var ScriptsCommand = &passport.Command{
	Name:        "scripts",
	Description: "provides commands used to manage workspace scripts",
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

		name := cmd.Args.String("name")
		if name == "" {
			cmd.Help()
			return nil
		}

		s, err := w.GetScript(name)
		if err != nil {
			return err
		}

		fmt.Printf("Command: %v\n", s.Command)

		return nil
	},
	Args: passport.CommandArgs{
		{
			Name:        "name",
			Description: "optionally, a script name can be passed to give info",
		},
	},
	Cmds: passport.CommandSet{
		listScriptsCommand,
		addScriptsCommand,
	},
}
