package passport

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandSet_ParseCommand(t *testing.T) {
	set := CommandSet{
		{
			Name: "version",
		},
		{
			Name: "math",
			Cmds: []*Command{
				{
					Name: "add",
				},
			},
		},
	}

	t.Run("Get Version Command", func(t *testing.T) {
		cmd := set.ParseCommand([]string{"version"})
		assert.Equal(t, set[0], cmd)
	})

	t.Run("Get Math Add Command", func(t *testing.T) {
		args := []string{"math", "add", "-x", "1", "-y", "2"}
		cmd := set.ParseCommand(args)
		assert.Equal(t, set[1].Cmds[0], cmd)
	})

	t.Run("Given Invalid Args", func(t *testing.T) {
		args := []string{"test", "command"}
		cmd := set.ParseCommand(args)
		assert.Nil(t, cmd)
	})

	t.Run("Given Nil Args", func(t *testing.T) {
		cmd := set.ParseCommand(nil)
		assert.Nil(t, cmd)
	})

	t.Run("Given 0 Args", func(t *testing.T) {
		cmd := set.ParseCommand([]string{})
		assert.Nil(t, cmd)
	})
}

func TestCommandSet_Help(t *testing.T) {
	cs := CommandSet{
		{
			Name: "TestCommand1",
		},
		{
			Name: "TestCommand2",
		},
	}

	pr, pw, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	oldStdout := os.Stdout
	os.Stdout = pw

	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	cs.Help()

	pw.Close()
	os.Stdout = oldStdout

	outputBytes, err := ioutil.ReadAll(pr)
	if err != nil {
		panic(err)
	}

	output := string(outputBytes)

	assert.Contains(t, output, "TestCommand1\n")
	assert.Contains(t, output, "TestCommand2\n")
}

func TestCommand_ParseArgs(t *testing.T) {
	nameArg := &CommandArg{Name: "name"}
	flagArg := &CommandArg{Name: "my-flag", IsFlag: true}

	cmd := &Command{
		Args: []*CommandArg{
			nameArg,
			flagArg,
		},
	}

	args := []string{"--my-flag", "--name", "reece"}
	cmd.ParseArgs(args)

	assert.Equal(t, "reece", nameArg.Value)
	assert.Equal(t, "true", flagArg.Value)
}

func TestCommand_Help(t *testing.T) {
	cmd := &Command{
		Name:        "TestCommand",
		Description: "MyTestCommand",
		Args: []*CommandArg{
			{
				Name:        "TestArg1",
				Description: "MyTestArg1",
				IsFlag:      true,
			},
		},
		Cmds: CommandSet{
			{
				Name:        "TestCommand1",
				Description: "MyTestCommand1",
				Args: []*CommandArg{
					{
						Name:        "TestArg1",
						Description: "MyTestArg1",
						IsFlag:      false,
					},
				},
			},
		},
	}

	pr, pw, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	oldStdout := os.Stdout
	os.Stdout = pw

	t.Cleanup(func() {
		os.Stdout = oldStdout
	})

	cmd.Help()

	pw.Close()
	os.Stdout = oldStdout

	outputBytes, err := ioutil.ReadAll(pr)
	if err != nil {
		panic(err)
	}

	output := string(outputBytes)

	assert.Contains(t, output, "TestCommand\n")
	assert.Contains(t, output, "---\nMyTestCommand\n")
	assert.Contains(t, output, "\targs:\n")
	assert.Contains(t, output, "TestArg1:\tMyTestArg1\n")
	assert.Contains(t, output, "flag: true\n")
	assert.Contains(t, output, "commands:\n")
	assert.Contains(t, output, "TestCommand1:\tMyTestCommand1\n")
	assert.Contains(t, output, "\t\t\targs:\n")
	assert.Contains(t, output, "TestArg1:\tMyTestArg1\n")
	assert.Contains(t, output, "flag: false\n")
}

func TestCommandArgs_String(t *testing.T) {
	args := CommandArgs{
		{
			Name:  "test",
			Value: "Hello World",
		},
	}

	t.Run("Given Valid Arg Name", func(t *testing.T) {
		v := args.String("test")
		assert.Equal(t, "Hello World", v)
	})

	t.Run("Where Arg Does Not Exist", func(t *testing.T) {
		v := args.String("hello-world")
		assert.Equal(t, "", v)
	})
}

func TestCommandArgs_Bool(t *testing.T) {
	args := CommandArgs{
		{
			Name:  "test1",
			Value: "true",
		},
		{
			Name:  "test2",
			Value: "false",
		},
		{
			Name:  "test3",
			Value: "not a bool value",
		},
	}

	t.Run("Where Value Should Be True", func(t *testing.T) {
		v := args.Bool("test1")
		assert.True(t, v)
	})

	t.Run("Where Value Should Be False", func(t *testing.T) {
		v := args.Bool("test2")
		assert.False(t, v)
	})

	t.Run("Where Value Is Not Bool", func(t *testing.T) {
		v := args.Bool("test3")
		assert.False(t, v)
	})

	t.Run("Where Arg Does Not Exist", func(t *testing.T) {
		v := args.Bool("hello-world")
		assert.False(t, v)
	})
}
