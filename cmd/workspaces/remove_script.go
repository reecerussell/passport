package workspaces

import (
	"fmt"
	"os"

	"github.com/reecerussell/passport"
)

var removeScriptsCommand = &passport.Command{
	Name:        "rm",
	Description: "used remove a new script from a workspace",
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

		err = w.RemoveScript(name)
		if err != nil {
			return err
		}

		err = cnf.Save()
		if err != nil {
			return err
		}

		fmt.Println("Successfully removed the script!")

		return nil
	},
	Args: passport.CommandArgs{
		{
			Name:        "name",
			Description: "the name of the script to remove",
			IsFlag:      false,
		},
	},
}
