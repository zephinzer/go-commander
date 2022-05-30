package commander

import (
	"io"
)

// NewCommand initialises an instance of the command
// class and returns the Command interface which
// prevents accidental modification of the underlying
// data
//
// Example for `ls -al` follows:
//
//   commander.NewCommand("ls").
//     AddParams("-al").
//     Execute()
//
// Example for `du -d 1 -h` follows:
//
//   commander.NewCommand("du").
//     AddParams("-d", "1").
//     AddParams("-h").
//     Execute()
func NewCommand(invocation string) Command {
	cmd := command{
		Environment: map[string]string{},
		Invocation:  invocation,
		Stdout:      []io.Writer{},
		Stderr:      []io.Writer{},
	}
	return &cmd
}
