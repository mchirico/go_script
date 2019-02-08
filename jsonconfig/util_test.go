package jsonconfig

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

var tmpFile = "./tmpFile"

func TestMain(m *testing.M) {
	tmpFile := SetupFunction()
	code := m.Run()
	os.Exit(code)
	TeardownFunction(tmpFile)
}

func SetupFunction() string {

	input := `
        {
            "squadName": "Super hero squad",
            "homeTown": "Metro City",
            "formed": 2016,
            "secretBase": "Super tower",
            "active": true
}`

	WriteFile(input, tmpFile)
	return tmpFile
}

func TeardownFunction(tmpFile string) {
	os.RemoveAll(tmpFile)

}

func TestReadFileError(t *testing.T) {

	var str bytes.Buffer
	log.SetOutput(&str)

	file := "///"
	m := map[string]string{}
	err := ReadJSON(file, m)
	if err == nil {
		t.FailNow()
	}

	expected := `Can't read file`

	if strings.Contains(str.String(), expected) != true {
		t.Logf("Expected: %s\n", expected)
		t.Logf("Got: %s\n", str.String())
	}

}

func TestReadJSONError(t *testing.T) {

	var str bytes.Buffer
	log.SetOutput(&str)

	file := "_junk_bad"
	d1 := []byte("hello\ngo\n")
	ioutil.WriteFile(file, d1, 0644)

	m := map[string]string{}
	err := ReadJSON(file, m)
	if err == nil {
		t.FailNow()
	}

	expected := `invalid character 'h' looking for beginning of value`
	if strings.Contains(err.Error(), expected) != true {
		t.Logf("Expected: %s\n", expected)
		t.Logf("Got: %s\n", err.Error())
		t.FailNow()
	}

	expected = `json.Unmarshal error`
	if strings.Contains(str.String(), expected) != true {
		t.Logf("Expected: %s\n", expected)
		t.Logf("Got: %s\n", str.String())
	}

}

func TestReadJSON(t *testing.T) {
	type squad struct {
		SquadName  string
		HomeTown   string
		Formed     int
		SecretBase string
		Active     bool
	}

	expected := "Super hero squad"

	s := squad{}
	ReadJSON(tmpFile, &s)
	if s.SquadName != expected {
		t.Fatalf("Expected:\n%s\n\nGot:\n%s\n", expected, s.SquadName)
	}
	log.Printf("%v\n", s.SquadName)
}
