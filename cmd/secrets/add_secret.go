package secrets

import (
	"github.com/reecerussell/passport"
)

var addSecretCommand = &passport.Command{
	Name:        "add",
	Description: "used to add a secret",
	Execute: func(cmd *passport.Command, ctx *passport.CommandContext) error {
		cnf, err := passport.LoadConfig(ctx.ConfigDir, ctx.Fs)
		if err != nil {
			return err
		}

		name := cmd.Args.String("name")
		value := cmd.Args.String("value")
		plainText := cmd.Args.Bool("plain-text")

		err = cnf.AddSecret(name, value, !plainText, ctx.Crypto)
		if err != nil {
			return err
		}

		return cnf.Save()
	},
	Args: passport.CommandArgs{
		{
			Name:        "name",
			Description: "the name of the secret",
		},
		{
			Name:        "value",
			Description: "the value of the secret",
		},
		{
			Name:        "plain-text",
			Description: "determines whether the value should be stored in plain text",
			IsFlag:      true,
		},
	},
}
