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
	t.Fatalf("❌ Error: process ran with err %v, want exit status 1", err)
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

func TestScript_LogProcess_TimeDelta(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(),
		100000*time.Millisecond)
	defer cancel()

	s := Script{}
	s.JSON.Command = `sleep 3`
	s.JSON.Log = tmpFile
	s.LogProcess(ctx)

	expected := int64(13)
	if s.DeltaTime.FileSize != expected {
		t.Fatalf("❌ Expected filesize: %d Got: %d\n", expected, s.DeltaTime.FileSize)
	}

	if s.DeltaTime.D1 >= 3*time.Second && s.DeltaTime.D1 <= 4*time.Second {
		t.Logf("Time good")
	} else {
		t.FailNow()
	}
	t.Logf("✅ Successful execution of TestScript_LogProcess_TimeDelta")
}

func TestZeroOut(t *testing.T) {

	if n, err := ZeroOut(tmpFile); err != nil && n != 0 {
		t.FailNow()
	}

	if n, tot, err := WriteData(tmpFile, []byte("This is a test")); err !=
		nil || tot <= 0 || n <= 0 {
		t.Errorf("❌ Can't write data to temp file")
	}

	if n, err := ZeroOut(tmpFile); err != nil && n != 0 {
		t.Fatalf("❌ Error")
	}

}

func TestLoopWithTimeout(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(),
		20*time.Millisecond)
	defer cancel()

	s := Script{}
	s.JSON.Command = `body() { IFS= read -r header; printf '%s %s\n %s\n' $(date "+%Y-%m %H:%M:%S") "$header"; "$@"; } && ps aux| body sort -n -r -k 4|head -n4`
	s.JSON.Log = tmpFile
	s.JSON.LogSizeLimit = 20000

	s.Loop(ctx, 1000)

	data, err := ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("❌ Got err reading %v err: %v\n", tmpFile, err)
	}

	if strings.Contains(data, "STAT START") != true {
		t.Fatalf("❌ Expected: %v, got: %v\n", "STAT STARTED", data)
	}
	t.Logf("✅ Successful execution of TestLoopWithTimeout")
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
		t.Fatalf("❌ Can't work working directory")
	}
}

func TestSpaceAvailable(t *testing.T) {
	space := SpaceAvailable(".")
	if space <= 0 {
		t.Fatalf("❌ Error: TestSpaceAvailable")
	}
	t.Logf("✅ Successful execution of TestSpaceAvailable")
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
		t.Fatalf("❌ Error")
	}

	if strings.Contains(string(dat), "BEGIN") != true {
		t.Logf("BEGIN not found. Got: %s", string(dat))
		t.Fatalf("❌ Error")
	}

	if strings.Contains(str.String(), "fileSize: 10076") != true {
		t.Fatalf("❌ Error: %s", str.String())
	}

	t.Logf("✅ Successful execution of TestWriteData_Append")
}

// Test fails in travis
func TestWriteData_Error(t *testing.T) {
	var str bytes.Buffer

	log.SetOutput(&str)

	a, b, err := WriteData("/DummyDIR/ShouldFail//", []byte("junk"))
	if a != -1 && b != -1 && err == nil {
		t.Fatalf("❌ TRUE failure..")
	}

	if strings.Contains(str.String(), "Error WriteData. os.OpenFile open") != true {
		t.Fatalf("❌ Expected:\n%s\n\nGot:\n%s\n",
			"Error WriteData. os.OpenFile open", str.String())
	}

	t.Logf("✅ TestWriteData_Error SUCCESS!!")
}

func TestLast(t *testing.T) {

	t.Logf("last test..")
	t.Logf("Trying to find out why last test some times")
	t.Logf("fails in travis.")
}
