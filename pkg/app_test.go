package pkg

import (
	"bytes"
	"context"
	"fmt"
	"github.com/mchirico/go_script/analyze"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"
)

var tmpFile = "/tmp/p.log"
var dir = "./tmp"

func TestMain(m *testing.M) {
	SetupFunction()
	retCode := m.Run()
	TeardownFunction()
	os.Exit(retCode)
}

func SetupFunction() {
	//MakeDir(dir)
}

func TeardownFunction() {
	os.RemoveAll(tmpFile)
	//os.RemoveAll(dir)
}

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
	s.JSON.Command = `sleep 3`
	s.JSON.Log = tmpFile
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

func TestLoopWithTimeout(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(),
		20*time.Millisecond)
	defer cancel()

	s := Script{}
	s.JSON.Command = `body() { IFS= read -r header; printf '%s %s\n %s\n' $(date "+%Y-%m %H:%M:%S") "$header"; "$@"; } && ps aux| body sort -n -r -k 4|head -n4`
	s.JSON.Log = tmpFile
	s.JSON.LogSizeLimit = 2000

	s.Loop(ctx, 1000)

	data, err := ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Got err reading %v err: %v\n", tmpFile, err)
	}

	if strings.Contains(data, "STAT START") != true {
		t.Fatalf("Expected: %v, got: %v\n", "STAT STARTED", data)
	}

}

// TODO: Implement this...
func TestLoopSize(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(),
		2000*time.Millisecond)
	defer cancel()

	s := Script{}

	fmt.Printf("s.Analyze=%v", s.Analyze)
	s.Analyze = analyze.Print
	s.JSON.Command = `ps aux|head -n4`
	s.JSON.LoopDelay = 2
	s.JSON.Log = tmpFile

	s.Loop(ctx, 1000)

}

func TestGetDir(t *testing.T) {
	dir := GetDir("junk")
	if strings.Contains(dir, "/") != true {
		t.Fatalf("Can't work working directory")
	}
}

func TestSpaceAvailable(t *testing.T) {
	space := SpaceAvailable(".")
	if space <= 0 {
		t.FailNow()
	}

}

func TestWriteLog(t *testing.T) {
	var str bytes.Buffer

	log.SetOutput(&str)
	log.Print("test")

	fmt.Printf("Here's the log message: '%v'\n",
		strings.TrimSuffix(str.String(), "\n"))
}

func TestWriteData_Append(t *testing.T) {

	var str bytes.Buffer
	log.SetOutput(&str)

	file := "tmpFile"
	os.RemoveAll(file)

	data := []byte("BEGIN\n")
	WriteData(file, data)

	for i := 0; i < 300; i++ {
		data = []byte(fmt.Sprintf("Filler Stuff: %d\n", i))
		WriteData(file, data)
	}

	data = []byte("Final\n")
	WriteData(file, data)

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		t.FailNow()
	}

	if strings.Contains(string(dat), "BEGIN") != true {
		t.Logf("BEGIN not found. Got: %s", string(dat))
		t.FailNow()
	}

	if strings.Contains(str.String(), "fileSize: 10076") != true {
		t.Logf("%s", str.String())
		t.FailNow()
	}

}

func TestWriteData_Error(t *testing.T) {
	var str bytes.Buffer

	log.SetOutput(&str)

	a, b, err := WriteData("////", []byte("junk"))
	if a != -1 && b != -1 && err == nil {
		t.FailNow()
	}

	if strings.Contains(str.String(), "Error WriteData. os.OpenFile open") != true {
		t.Logf("Expected:\n%s\n\nGot:\n%s\n",
			"Error WriteData. os.OpenFile open", str.String())
		t.FailNow()
	}

}
