package workspaces

import (
	"fmt"
	"os"

	"github.com/reecerussell/passport"
)

var listScriptsCommand = &passport.Command{
	Name:        "ls",
	Description: "used to list scripts in a workspace",
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

		fmt.Printf("Workspace: %s\n", w.Name)
		fmt.Println("Scripts:")

		for _, s := range w.Scripts {
			fmt.Printf("> %s\n", s.Name)
		}

		return nil
	},
}
