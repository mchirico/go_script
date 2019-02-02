package pkg

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"
)

var tmpFile = "/tmp/p.log"

func TestLogProcessKill(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		logProcessTimeout()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestLogProcessKill")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func logProcessTimeout() {

	ctx, cancel := context.WithTimeout(context.Background(),
		10*time.Millisecond)
	defer cancel()

	s := Script{}
	s.Command = `body() { IFS= read -r header; printf '%s %s\n %s\n' $(date "+%Y-%m %H:%M:%S") "$header"; "$@"; } && ps aux| body sort -n -r -k 4|head -n4`
	s.Log = tmpFile
	s.LogProcess(ctx)

}

func TestZeroOut(t *testing.T) {

	if n, err := ZeroOut(tmpFile); err != nil && n != 0 {
		t.FailNow()
	}

	if n, tot, err := WriteData(tmpFile, []byte("This is a test")); err !=
		nil || tot <= 0 || n <= 0 {
		t.Errorf("Can't write data to temp file")
	}

	if n, err := ZeroOut(tmpFile); err != nil && n != 0 {
		t.FailNow()
	}

}
