package passport

import "fmt"

// ExecuteFunc is a function type used to define how a Command is executed.
type ExecuteFunc func(cmd *Command, ctx *CommandContext) error

// Command is a task that can be executed. A command contains
// information about the command, in addition to an execute function,
// command arguments and sub-commands.
type Command struct {
	Name        string
	Description string
	Execute     ExecuteFunc
	Args        CommandArgs
	Cmds        CommandSet
}

// ParseArgs deserialises args into the command's arguments.
func (cmd *Command) ParseArgs(args []string) {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg[0] != '-' {
			continue
		}

		for _, cmdArg := range cmd.Args {
			if "--"+cmdArg.Name != arg {
				continue
			}

			if len(args) > i+1 && args[i+1][0] != '-' && !cmdArg.IsFlag {
				cmdArg.Value = args[i+1]
				i++
			}

			if cmdArg.IsFlag {
				cmdArg.Value = "true"
			}
		}
	}
}

// Help prints information about the command to os.Stdout.
func (cmd *Command) Help() {
	fmt.Printf("%s\n", cmd.Name)
	fmt.Println("---")
	fmt.Println(cmd.Description)

	if len(cmd.Args) > 0 {
		fmt.Printf("\targs:\n")

		for _, arg := range cmd.Args {
			fmt.Printf("\t\t%s:\t%s\n", arg.Name, arg.Description)
			fmt.Printf("\t\t\tflag: %v\n", arg.IsFlag)
		}
	}

	if len(cmd.Cmds) > 0 {
		fmt.Printf("\tcommands:\n")

		for _, sc := range cmd.Cmds {
			fmt.Printf("\t\t%s:\t%s\n", sc.Name, sc.Description)

			if len(sc.Args) > 0 {
				fmt.Printf("\t\t\targs:\n")

				for _, arg := range sc.Args {
					fmt.Printf("\t\t\t\t%s:\t%s\n", arg.Name, arg.Description)
					fmt.Printf("\t\t\t\t\tflag: %v\n", arg.IsFlag)
				}
			}
		}
	}
}

// CommandArg represents a command line argument, with relevent
// information for documentation and deserialization.
type CommandArg struct {
	Name        string
	Description string
	Value       string
	IsFlag      bool
}

// CommandContext is a struct provided to each command's execute
// function, providing common values.
type CommandContext struct {
	ConfigDir string
	Crypto    CryptoProvider
}

// CommandSet is a wrapper around []*Command, which provides helper functions.
type CommandSet []*Command

// ParseCommand returns a command matching the given args. If a command
// is found, the args are then parsed and the command is returned.
// Otherwise, nil is returned.
func (set CommandSet) ParseCommand(args []string) *Command {
	if args == nil || len(args) < 1 {
		return nil
	}

	arg := args[0]

	for _, cmd := range set {
		if cmd.Name != arg {
			continue
		}

		if len(args) > 1 && args[1][0] != '-' {
			c := cmd.Cmds.ParseCommand(args[1:])
			if c != nil {
				c.ParseArgs(args[2:])
				return c
			}
		}

		cmd.ParseArgs(args[1:])
		return cmd
	}

	return nil
}

// Help calls cmd.Help() on each of the commands in the set.
func (set CommandSet) Help() {
	for _, cmd := range set {
		cmd.Help()
		fmt.Printf("\n")
	}
}

// CommandArgs is a wrapper around []*CommandArgs, which provides helper functions.
type CommandArgs []*CommandArg

// String returns a string value for an argument with the given name.
// If the argument could not be found, an empty string is returned.
func (args CommandArgs) String(name string) string {
	for _, arg := range args {
		if arg.Name == name {
			return arg.Value
		}
	}

	return ""
}

// Bool returns a boolean value for an argument with the given name.
// If the argument could not be found, false is returned.
func (args CommandArgs) Bool(name string) bool {
	for _, arg := range args {
		if arg.Name == name {
			return arg.Value == "true"
		}
	}

	return false
}
