package jsonconfig

import (
	"log"
	"os"
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
