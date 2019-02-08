package main

import (
	"bytes"
	"encoding/json"
	"github.com/mchirico/go_script/jsonconfig"
	"github.com/mchirico/go_script/pkg"
	"github.com/mchirico/go_script/yamlpkg"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	//a = pkg.App{}

	//a.Initilize()
	code := m.Run()

	os.Exit(code)
}

var file = "_fixture_test.yaml"

// There's an issue with format... going back and forth
func TestReadJson(t *testing.T) {

	s := pkg.Script{}
	j := s.JSON

	j.LoopDelay = 20
	j.Log = "/tmp/mem.log"
	j.LogSizeLimit = 40000
	j.Command = `./cmd.sh`
	j.Command = `body() { IFS= read -r header; printf '%s %s\n %s\n' $(date "+%Y-%m %H:%M:%S") "$header"; "$@"; } && ps aux| body sort -n -r -k 4|head -n4`

	j.ArchiveLog = "/tmp/archive.log"
	j.DieAfterHours = 200

	odata, err := json.Marshal(j)
	if err != nil {
		log.Printf("ejson: %s\n", err)
	}

	file := ".script"

	log.Println(string(odata))

	n, err := jsonconfig.WriteFile(string(odata),
		file)
	if err != nil {
		log.Printf("error: %v, %v\n", n, err)
		t.Fail()
	}

	j = pkg.JSON{}
	err = jsonconfig.ReadJSON(file, &j)
	if err != nil {
		t.Fatalf("s=%v", s)
	}
}

func Test_CreateS(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		os.RemoveAll(file)
		createS(file)
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=Test_CreateS")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {

		// Now go though 2nd run

		expected := `body() { IFS= read -r header; printf '%s %s\n %s\n' $(date "+%Y-%m %H:%M:%S") "$header"; "$@"; } && ps aux| body sort -n -r -k 4`

		s := createS(file)
		if strings.Contains(s.JSON.Command, expected) != true {
			t.Fatalf("Expected:\n%s\n\nGot:\n%s\n\n",
				s.JSON.Command, expected)
		}
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func Test_Main(t *testing.T) {

	var str bytes.Buffer
	log.SetOutput(&str)

	c := yamlpkg.Config{}
	c.SetDefault()
	c.Write(configFile)

	err := c.Read(configFile)
	if err != nil {
		t.FailNow()
	}

	c.Yaml.DieAfterHours = 0
	c.Yaml.LoopDelay = 1
	c.Yaml.LogSizeLimit = 300
	c.Write(configFile)
	main()

	expected := []string{}
	expected = append(expected, "Zero out called:")
	expected = append(expected, "Space Available")
	expected = append(expected, "wrote:")
	expected = append(expected, "fileSize:")

	for index, _ := range expected {
		if strings.Contains(str.String(), expected[index]) != true {
			t.FailNow()
		}
	}

}
