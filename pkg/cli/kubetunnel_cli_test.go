package cli_test

import (
	"bytes"
	"testing"
)

func TestExecuteHelp(t *testing.T) {

	//cmd := cli.NewRootCmd()
	stdout := bytes.NewBufferString("")
	_ = stdout
	//cmd.SetOut(stdout)
	//
	//cmd.SetArgs(strings.Split("kubetunnel --help", " "))
	//
	//cmd.Execute()
	//
	//out, err := ioutil.ReadAll(stdout)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//if string(out) != "hi-via-args" {
	//	t.Fatalf("expected \"%s\" got \"%s\"", "hi-via-args", string(out))
	//}
}
