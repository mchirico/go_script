package yamlpkg

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestReadWrite(t *testing.T) {
	file := "script.yaml"

	c := Config{}

	c.Yaml.Command = `body() { IFS= read -r header; printf '%s %s\n %s\n' $(date "+%Y-%m %H:%M:%S") "$header"; "$@"; } && ps aux| body sort -n -r -k 4|head -n4`
	c.Yaml.Log = "mem.log"
	c.Yaml.LoopDelay = 20
	c.Yaml.LogSizeLimit = 40000
	c.Yaml.ArchiveLog = "memarchive.log"
	c.Yaml.DieAfterHours = 200

	c.Write(file)

	c2 := Config{}
	c2.Read(file)

	if c2.Yaml.Command != c.Yaml.Command {
		t.Fatalf("Got:\n%s\n\nExpected:\n%s\n", c2.Yaml.Command, c.Yaml.Command)
	}

}

func TestConfig_SetDefault(t *testing.T) {
	c := Config{}
	c.SetDefault()
	if c.Yaml.LoopDelay != 20 {
		t.FailNow()
	}

}
