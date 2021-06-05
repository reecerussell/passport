package secrets

import (
	"github.com/reecerussell/passport"
)

var removeSecretCommand = &passport.Command{
	Name:        "rm",
	Description: "used to remove a secret",
	Execute: func(cmd *passport.Command, ctx *passport.CommandContext) error {
		cnf, err := passport.LoadConfig(ctx.ConfigDir)
		if err != nil {
			return err
		}

		name := cmd.Args.String("name")
		err = cnf.RemoveSecret(name)
		if err != nil {
			return err
		}

		return cnf.Save()
	},
	Args: passport.CommandArgs{
		{
			Name:        "name",
			Description: "the name of the secret to remove",
		},
	},
}
