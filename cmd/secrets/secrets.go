package secrets

import (
	"fmt"

	"github.com/reecerussell/passport"
)

// Command is the main entrypoint command for operations around secrets.
var Command = &passport.Command{
	Name:        "secrets",
	Description: "provides commands used to manage and view secrets",
	Execute: func(cmd *passport.Command, ctx *passport.CommandContext) error {
		cnf, err := passport.LoadConfig(ctx.ConfigDir)
		if err != nil {
			return err
		}

		name := cmd.Args.String("name")
		if name == "" {
			cmd.Help()
			return nil
		}

		s, err := cnf.GetSecret(name)
		if err != nil {
			return err
		}

		fmt.Printf("Name: %s\n", s.Name)
		fmt.Printf("Value: %s\n", s.GetValue())
		fmt.Printf("Secure: %v\n", s.Secure)

		return nil
	},
	Args: passport.CommandArgs{
		{
			Name:        "name",
			Description: "optionally, a name can be passed in to get a single secret",
		},
	},
	Cmds: passport.CommandSet{
		listSecretsCommand,
		addSecretCommand,
		removeSecretCommand,
	},
}
