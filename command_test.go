package passport

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseCommand(t *testing.T) {
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
}

func TestParseArgs(t *testing.T) {
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
