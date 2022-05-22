package commander

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestCommand(t *testing.T) {
	suite.Run(t, &CommandTests{})
}

type CommandTests struct {
	suite.Suite
}

func (s CommandTests) Test_AddParam() {
	cmd := NewCommand("test").AddParam("a")
	cmdInstance, ok := cmd.(*command)
	s.True(ok)
	s.Equal("a", cmdInstance.Params[0].Label)
	s.Nil(cmdInstance.Params[0].Value)
}

func (s CommandTests) Test_AddParam_withValue() {
	cmd := NewCommand("test").AddParam("a", "b")
	cmdInstance, ok := cmd.(*command)
	s.True(ok)
	s.Equal("b", *cmdInstance.Params[0].Value)
}

func (s CommandTests) Test_AddStderr() {
	altStderrBuffer := bytes.NewBuffer(nil)
	altStderr := bufio.NewWriter(altStderrBuffer)
	cmd := NewCommand("test").AddStderr(altStderr)
	cmdInstance, ok := cmd.(*command)
	s.True(ok)
	s.Len(cmdInstance.Stderr, 1)
}

func (s CommandTests) Test_AddStdout() {
	altStdoutBuffer := bytes.NewBuffer(nil)
	altStdout := bufio.NewWriter(altStdoutBuffer)
	cmd := NewCommand("test").AddStdout(altStdout)
	cmdInstance, ok := cmd.(*command)
	s.True(ok)
	s.Len(cmdInstance.Stdout, 1)
}

func (s CommandTests) Test_GetAsString() {
	commandAsString :=
		NewCommand("test").
			AddParam("--label1").
			AddParam("--label2", "value2").
			AddParam("--label3", `"value3"`).
			AddParam("--label4", "value 4").
			AddParam("--label5", `"value 5"`).
			GetAsString()
	s.EqualValues(strings.Trim(`
test \
  --label1 \
  --label2 value2 \
  --label3 "value3" \
  --label4 "value 4" \
  --label5 ""value 5"";
`, "\n"), commandAsString)
}

func (s CommandTests) Test_GetAsString_oneliner() {
	commandAsString :=
		NewCommand("test").
			AddParam("--label1").
			AddParam("--label2", "value2").
			AddParam("--label3", `"value3"`).
			AddParam("--label4", "value 4").
			AddParam("--label5", `"value 5"`).
			GetAsString(true)
	s.EqualValues(`test --label1 --label2 value2 --label3 "value3" --label4 "value 4" --label5 ""value 5"";`, commandAsString)
}

func (s CommandTests) Test_Run() {
	cmd := NewCommand("ls").
		AddParam("-ls").
		AddParam("./tests/lsal")
	cmdOutput := cmd.Execute()
	s.Nil(cmdOutput.Error)
	outputValue := cmdOutput.Stdout.String()
	s.EqualValues(strings.TrimLeft(`
file1
file2
folder_a
folder_b
`, "\n"), outputValue)
}
