package commander

import (
	"bytes"
	"io"
)

type Command interface {
	// AddParam adds a new parameter to the command,
	// each parameter can either be a bool flag
	// (eg. --enable-something) or a value flag
	// (eg. --some-value 1). Specify the :value
	// argument to specify a value flag
	AddParam(key string, value ...string) Command

	// AddStderr adds an additional buffer stream to
	// send standard error to
	AddStderr(output io.Writer) Command

	// AddStdout adds an additional buffer stream to
	// send standard output to
	AddStdout(output io.Writer) Command

	// DisableGlobalEnvironment when specified
	// disables the global environment from being
	// injected into the command's environment. You
	// would typically want to inject the environment
	// therefore this is implemented semantically as
	// a DisableX instead of an EnabledX which is more
	// intuitive to humans generally speaking
	DisableGlobalEnvironment() Command

	// EnableStderr enables printing direct to os.Stderr
	EnableStderr() Command

	// EnableStdout enables printing direct to os.Stdout
	EnableStdout() Command

	// Execute runs the command
	Execute() output

	// GetAsString retrieves the command in string
	// form, can be used for things like showing the
	// command to be run to the user or for saving the
	// comamnd into some file
	GetAsString(oneLine ...bool) string

	// SetEnvironment sets the environment variable
	// identified by the :key to the value :value
	SetEnvironment(key string, value string) Command

	// SetStderr sets the streams for standard error
	// to be sent to
	SetStderr(output io.Writer) Command

	// SetStdout sets the streams for standard output
	// to be sent to
	SetStdout(output io.Writer) Command

	// SetWorkingDirectory sets the working directory
	// for the command to be run in
	SetWorkingDirectory(wd string) Command
}

type command struct {
	// Environment if specified will be injected to run
	// the command. If UseGlobalEnvironment is set to
	// truthy, this set of environment variables will be
	// injected after the global environment is injected
	Environment map[string]string

	// Invocation is a command that should be available
	// in the user's $PATH
	Invocation string

	// isGlobalEnvironmentDisabled will inject all environment
	// variables from the user's system to this command
	// when it's being run. Environment values from the
	// Environment property will take precedence
	isGlobalEnvironmentDisabled bool

	// isStderrEnabled when set to true disables stderr
	// from being streamed to the console
	isStderrEnabled bool

	// isStdoutEnabled when set to true disables stdout
	// from being streamed to the console
	isStdoutEnabled bool

	// Params are parameters to the invocation
	Params []Param

	// Stdout is a list of extra io.Writer interfaces
	// where output to stdout will additionally be sent to
	Stdout []io.Writer

	// Stderr is a list of extra io.Writer interfaces
	// where output to stderr will additionally be sent to
	Stderr []io.Writer

	// WorkingDirectory specifies where the command will
	// be executed from
	WorkingDirectory string
}

type Param struct {
	// Label is the --label porttion of a flag
	Label string

	// Value is an optional value that comes after the
	// flag lebel is specified (eg. --label value)
	Value *string
}

type output struct {
	Error  error
	Stdout *bytes.Buffer
	Stderr *bytes.Buffer
}
