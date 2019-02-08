/*

useradd -r -M -s /bin/false analyze-goguy
mkdir -p /webproject/analyze-goguy
*/

package main

import (
	"context"
	"flag"
	"github.com/mchirico/go_script/pkg"
	"github.com/mchirico/go_script/yamlpkg"
	"log"
	"time"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "configFile", "script.yaml", "Yaml Configuration file")
	flag.Parse()
}

func createS(configFile string) (pkg.Script, error) {
	s := pkg.Script{}

	c := yamlpkg.Config{}
	err := c.Read(configFile)
	if err != nil {
		c.SetDefault()
		c.Write(configFile)
		msg := `
Could not read script.yaml. Creating default.
You can run this command again to pickup default script.yaml

`
		log.Print(msg)
		return s, err
	}

	s.JSON.DieAfterHours = c.Yaml.DieAfterHours
	s.JSON.DieAfterSeconds = c.Yaml.DieAfterSeconds
	s.JSON.LogSizeLimit = c.Yaml.LogSizeLimit
	s.JSON.Log = c.Yaml.Log
	s.JSON.ArchiveLog = c.Yaml.ArchiveLog
	s.JSON.LoopDelay = c.Yaml.LoopDelay
	s.JSON.Command = c.Yaml.Command

	return s, err
}

func main() {

	s, err := createS(configFile)
	if err != nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(s.JSON.DieAfterHours)*time.Hour+time.Duration(3)*time.Second)
	defer cancel()

	s.Loop(ctx, 100000)

}
