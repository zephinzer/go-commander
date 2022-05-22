# Commander

This library wraps the `exec.Cmd` for readability and usability purposes. Some use cases I had when developing this:

1. I want to execute a shell command from a Go application
2. I want to get the output of the shell command to do some data transformations on it (eg. with the `aws` CLI command tool to list out infrastructure or with `kubectl` to get a list of pods/deployments *et cetera*)
3. I want to be able to inject the user's current environment into the command
4. I want the option to decide if the output to be printed to the default `stdout` or `stderr`

The purpose of this package is as an underlying library for a CLI tool that codifies contextual DevOps knowledge in an organisation. This means things like what commands should be run and in what order to achieve a certain desired outcome.

Since this can vary from organisation to organisation, it's usually not easy to pick up on the conventions and it's not useful to internalise this since different projects can have different conventions. Therefore, I write a CLI tooling around conventions, and this library helps with that.

# Usage

To use this library, import it:

```go
import "github.com/zephinzer/go-commander"
```

## Basic/universal usage

Create a new command via the `NewCommand(...)` function:

```go
func main() {
  // ...
  command := commander.NewCommand("ls").
    AddParam("-al")
  // ...
}
```

Finally, whenever you want, you can execute it:

```go
// ...
  output := command.Execute()
// ...
```

To get the output of the command, you can use the `.Stderr` and `.Stdout` properties of the output object:

```go
// ...
  // to get the stderr
  fmt.Println(output.Stderr.String())
  // to get the stdout
  fmt.Println(output.Stdout.String())
// ...
```

## Setting custom environment variables

It's usually useful to be able to inject environment variables into a custom command you wish to run. For example, if the user's `AWS_PROFILE` is set to `"production"` but you want it set to `"staging"`, you can inject environment variables using the `SetEnvironment` method. These will overwrite the global environment if the global environment is not disabled.

```go
// ...
  command := commander.NewCommand("aws").
    SetEnvironment("AWS_PROFILE", "staging").
    AddParam("sts")
    AddParam("get-caller-identity")
// ...
```

## Executing from a different directory

By default, the current working directory is used. If you would like to run the command from a different directory, use the `SetWorkingDirectory` method to set the path to the directory where you would like the command to run from

```go
// ...
  command := commander.NewCommand("ls").
    SetWorkingDirectory("./tests")
// ...
```

## Printing the command

For people who don't know what they're going to run, it can be useful to see the command they're running. This also gives you a better idea on how to debug things if things go wrong. To print the command as a string, you can use the `GetAsString` method

```go
// ...
  command := commander.NewCommand("ls").
    SetWorkingDirectory("./tests")

  fmt.Println(command.GetAsString())
// ...
```

## Enable `stdout`/`stderr` mirroring

To print the `stdout`/`stderr` output as it comes in, you can use the `EnableStdout` or `EnableStderr` chainable methods

```go
// ...
  command := commander.NewCommand("ls").
    SetWorkingDirectory("./tests").
    EnableStdout(). // and/or
    EnableStderr()
// ...
```

## Setting custom `stdout`/`stderr` streams

Occassionally it could be useful to stream the standard output/error to another buffer. To do this, you can use the `SetStderr` or `SetStdout` methods:

```go
// ...
  otherBuffer := bytes.NewBuffer(nil)
  otherBufferWriter := bufio.NewWriter(otherBuffer)

  command := commander.NewCommand("ls").
    SetWorkingDirectory("./tests").
    SetStderr(otherBufferWriter). // and/or
    SetStdout(otherBufferWriter)
// ...
```

## Disable environment variable injection

Sometimes it's useful to disable the global environment from polluting a command's environment. To do that, use the `.DisableGlobalEnvironment` method.

```go
// ...
  command := commander.NewCommand("ls").
    DisableGlobalEnvironment()
// ...
```

# Contribution

1. Run `go mod tidy` or `go mod vendor` to bring in the dependencies
2. Run `go test ./...` to test this package
3. When pushed to Gitlab, this repository will trigger a pipeline to run tests, see the [`./.gitlab-ci.yml` file](./.gitlab-ci.yml) for more information
