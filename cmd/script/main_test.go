package main

import (
	"encoding/json"
	"github.com/mchirico/go_script/jsonconfig"
	"github.com/mchirico/go_script/pkg"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	//a = pkg.App{}

	//a.Initilize()
	code := m.Run()

	os.Exit(code)
}

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
