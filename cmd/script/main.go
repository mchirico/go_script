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

func createS(configFile string) pkg.Script {
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
		log.Fatalf(msg)
	}

	s.JSON.DieAfterHours = c.Yaml.DieAfterHours
	s.JSON.LogSizeLimit = c.Yaml.LogSizeLimit
	s.JSON.Log = c.Yaml.Log
	s.JSON.ArchiveLog = c.Yaml.ArchiveLog
	s.JSON.LoopDelay = c.Yaml.LoopDelay
	s.JSON.Command = c.Yaml.Command

	return s
}

func main() {

	s := createS(configFile)

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(s.JSON.DieAfterHours)*time.Hour)
	defer cancel()

	s.Loop(ctx, 100000)

}
