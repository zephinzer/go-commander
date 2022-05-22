package commander

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"strings"
)

func (c *command) AddParam(key string, value ...string) Command {
	param := Param{
		Label: key,
	}
	if len(value) > 0 {
		param.Value = &value[0]
	}
	c.Params = append(c.Params, param)
	return c
}

func (c *command) AddStderr(output io.Writer) Command {
	c.Stderr = append(c.Stderr, output)
	return c
}

func (c *command) AddStdout(output io.Writer) Command {
	c.Stdout = append(c.Stdout, output)
	return c
}

func (c *command) EnableStderr() Command {
	c.isStderrEnabled = true
	return c
}

func (c *command) EnableStdout() Command {
	c.isStdoutEnabled = true
	return c
}

func (c *command) Execute() output {
	if c.WorkingDirectory == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return output{Error: fmt.Errorf("failed to get current working directory: %s", err)}
		}
		c.WorkingDirectory = cwd
	} else if !path.IsAbs(c.WorkingDirectory) {
		cwd, err := os.Getwd()
		if err != nil {
			return output{Error: fmt.Errorf("failed to get current working directory: %s", err)}
		}
		c.WorkingDirectory = path.Join(cwd, c.WorkingDirectory)
	}
	dirInfo, err := os.Lstat(c.WorkingDirectory)
	if err != nil {
		return output{Error: fmt.Errorf("failed to open directory at working directory '%s'", c.WorkingDirectory)}
	}
	if !dirInfo.IsDir() {
		return output{Error: fmt.Errorf("failed to find a directory at working directory path '%s'", c.WorkingDirectory)}
	}
	cmd := exec.Cmd{Dir: c.WorkingDirectory}

	// is what we calling actually there?
	fullPath, err := exec.LookPath(c.Invocation)
	if err != nil {
		return output{Error: fmt.Errorf("failed to find invocation '%s' in $PATH: %s", c.Invocation, err)}
	}
	if strings.Contains(fullPath, "/") {
		if !path.IsAbs(fullPath) {
			fullPath = path.Join(c.WorkingDirectory, fullPath)
		}
	}
	cmd.Path = fullPath

	// set arguments
	for _, param := range c.Params {
		cmd.Args = append(cmd.Args, param.Label)
		if param.Value != nil {
			cmd.Args = append(cmd.Args, *param.Value)
		}
	}

	// set environment
	if !c.isGlobalEnvironmentDisabled {
		cmd.Env = append(cmd.Env, os.Environ()...)
	}
	for key, value := range c.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}

	// prepare for the execution
	outputInstance := output{}

	// setup standard input stream
	cmd.Stdin = os.Stdin

	// setup standard error streams
	if c.isStderrEnabled {
		c.Stderr = append([]io.Writer{os.Stderr}, c.Stderr...)
	}
	outputInstance.Stderr = bytes.NewBuffer(nil)
	internalStderr := bufio.NewWriter(outputInstance.Stderr)
	defer internalStderr.Flush()
	cmd.Stderr = io.MultiWriter(append(c.Stderr, internalStderr)...)

	// setup standard output streams
	if c.isStdoutEnabled {
		c.Stdout = append([]io.Writer{os.Stdout}, c.Stdout...)
	}
	outputInstance.Stdout = bytes.NewBuffer(nil)
	internalStdout := bufio.NewWriter(outputInstance.Stdout)
	defer internalStdout.Flush()
	cmd.Stdout = io.MultiWriter(append(c.Stdout, internalStdout)...)

	outputInstance.Error = cmd.Run()

	return outputInstance
}

func (c *command) GetAsString(oneLine ...bool) string {
	var output strings.Builder

	delimiter := " \\\n  "
	if (len(oneLine) > 0) && oneLine[0] {
		delimiter = " "
	}

	output.WriteString(c.Invocation)
	for _, flag := range c.Params {
		if flag.Value == nil {
			output.WriteString(fmt.Sprintf("%s%s", delimiter, flag.Label))
		} else {
			printedValue := *flag.Value
			if strings.Contains(printedValue, "\"") {
				strings.ReplaceAll(printedValue, "\"", "\\\"")
			}
			if strings.Contains(printedValue, " ") {
				printedValue = fmt.Sprintf("\"%s\"", printedValue)
			}
			output.WriteString(fmt.Sprintf("%s%s %s", delimiter, flag.Label, printedValue))
		}
	}

	return output.String() + ";"
}

func (c *command) SetEnvironment(key string, value string) Command {
	c.Environment[key] = value
	return c
}

func (c *command) DisableGlobalEnvironment() Command {
	c.isGlobalEnvironmentDisabled = true
	return c
}

func (c *command) SetStderr(output io.Writer) Command {
	c.Stderr = []io.Writer{output}
	return c
}

func (c *command) SetStdout(output io.Writer) Command {
	c.Stdout = []io.Writer{output}
	return c
}

func (c *command) SetWorkingDirectory(wd string) Command {
	c.WorkingDirectory = wd
	return c
}
