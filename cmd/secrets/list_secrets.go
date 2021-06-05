package secrets

import (
	"fmt"

	"github.com/reecerussell/passport"
)

var listSecretsCommand = &passport.Command{
	Name:        "ls",
	Description: "used to list all secrets",
	Execute: func(cmd *passport.Command, ctx *passport.CommandContext) error {
		cnf, err := passport.LoadConfig(ctx.ConfigDir)
		if err != nil {
			return err
		}

		fmt.Println("Secrets:")

		for _, s := range cnf.Secrets {
			fmt.Printf("> %s\n", s.Name)
		}

		return nil
	},
}
