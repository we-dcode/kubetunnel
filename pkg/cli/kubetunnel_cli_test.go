package cli

import (
	"bytes"
	"io/ioutil"
	"strings"
	"testing"
)

func Test_ExecuteHelp(t *testing.T) {

	cmd := NewRootCmd()
	stdout := bytes.NewBufferString("")
	cmd.SetOut(stdout)

	cmd.SetArgs(strings.Split("kubetunnel --help", " "))

	cmd.Execute()

	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "hi-via-args" {
		t.Fatalf("expected \"%s\" got \"%s\"", "hi-via-args", string(out))
	}
}