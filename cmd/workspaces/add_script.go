package workspaces

import (
	"fmt"
	"os"

	"github.com/reecerussell/passport"
)

var addScriptsCommand = &passport.Command{
	Name:        "add",
	Description: "used add a new script to a workspace",
	Execute: func(cmd *passport.Command, ctx *passport.CommandContext) error {
		cnf, err := passport.LoadConfig(ctx.ConfigDir, ctx.Fs)
		if err != nil {
			return err
		}

		wd, _ := os.Getwd()
		w, err := cnf.GetWorkspace(wd)
		if err != nil {
			cnf.AddWorkspace(wd, wd)
			w, _ = cnf.GetWorkspace(wd)
		}

		name := cmd.Args.String("name")
		command := cmd.Args.String("command")

		err = w.AddScript(name, command)
		if err != nil {
			return err
		}

		err = cnf.Save()
		if err != nil {
			return err
		}

		fmt.Println("Successfully added new script!")

		return nil
	},
	Args: passport.CommandArgs{
		{
			Name:        "name",
			Description: "the name of the new script",
			IsFlag:      false,
		},
		{
			Name:        "command",
			Description: "the command to execute",
			IsFlag:      false,
		},
	},
}
