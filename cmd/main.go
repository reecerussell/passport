package main

import (
	"fmt"
	"os"
	"path"

	"github.com/reecerussell/passport"
	"github.com/reecerussell/passport/cmd/secrets"
	"github.com/reecerussell/passport/cmd/workspaces"
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

	fs := passport.NewFilesys()
	err := fs.EnsureDirectory(configDir)
	if err != nil {
		panic(err)
	}

	err = passport.EnsureConfigFile(configDir, fs)
	if err != nil {
		panic(err)
	}

	sets := passport.CommandSet{
		secrets.Command,
		workspaces.ScriptsCommand,
		workspaces.RunScriptCommand,
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
		Fs:        passport.NewFilesys(),
	}

	err = cmd.Execute(cmd, ctx)
	if err != nil {
		panic(err)
	}
}
