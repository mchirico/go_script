package analyze

import (
	"bytes"
	"io"
	"os"
	"testing"
)

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

}

func TestPrint(t *testing.T) {

	old := os.Stdout // keep backup of the real stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	msg := []byte(`Sample $@`)
	n, err := Print(msg)
	if n != 0 || err != nil {
		t.Fatalf("Print test returned error")
	}

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	os.Stdout = old // restoring the real stdout
	out := <-outC

	if out != string(msg)+"\n" {
		t.Fatalf("Expected:\n->%s<-\n\nGot:\n->%s<-", msg, out)
	}

	// reading our temp stdout
	//fmt.Println("previous output:")
	//fmt.Print(out)

}
