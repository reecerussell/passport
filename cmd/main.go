package main

import (
	"fmt"
	"os"
	"path"

	"github.com/reecerussell/passport"
	"github.com/reecerussell/passport/cmd/secrets"
)

var (
	version   string
	configDir string
)

func main() {
	userConfigDir, _ := os.UserConfigDir()
	if configDir == "" {
		configDir = path.Join(userConfigDir, ".passport")
	}

	if len(os.Args) <= 1 {
		fmt.Printf("Passport, %s\n", version)
		os.Exit(0)
	}

	err := passport.EnsureDirectory(configDir)
	if err != nil {
		panic(err)
	}

	err = passport.EnsureConfigFile(configDir)
	if err != nil {
		panic(err)
	}

	sets := passport.CommandSet{
		secrets.Command,
	}

	cmd := sets.ParseCommand(os.Args[1:])
	if cmd == nil {
		fmt.Println("Command not found!")
		sets.Help()
		os.Exit(1)
	}

	ctx := &passport.CommandContext{
		ConfigDir: configDir,
		Crypto:    passport.NewCryptoProvider(),
	}

	err = cmd.Execute(cmd, ctx)
	if err != nil {
		panic(err)
	}
}
