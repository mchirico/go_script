/*

useradd -r -M -s /bin/false analyze-goguy
mkdir -p /webproject/analyze-goguy
*/

package main

import (
	"context"
	"flag"
	"github.com/mchirico/go_script/jsonconfig"
	"github.com/mchirico/go_script/pkg"
	"log"
	"time"
)

func main() {

	var configFile string
	flag.StringVar(&configFile, "configFile", ".script", "Configuration file")
	flag.Parse()

	s := pkg.Script{}

	err := jsonconfig.ReadJSON(configFile, &s.JSON)
	if err != nil {

		log.Fatalf("cannot read config:")
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(s.JSON.DieAfterHours)*time.Hour)
	defer cancel()

	s.Loop(ctx, 1000, 20000)

}
