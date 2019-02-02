/*

useradd -r -M -s /bin/false analyze-goguy
mkdir -p /webproject/analyze-goguy
*/

package main

import (
	"context"
	"github.com/mchirico/go_script/pkg"
	"time"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(),
		1000*time.Millisecond)
	defer cancel()

	s := pkg.Script{}
	s.Command = `body() { IFS= read -r header; printf '%s %s\n %s\n' $(date "+%Y-%m %H:%M:%S") "$header"; "$@"; } && ps aux| body sort -n -r -k 4|head -n4`
	s.Log = "/tmp/p.log"

	for {
		_, size := s.LogProcess(ctx)

		time.Sleep(3 * time.Second)
		if size > 1533791 {
			pkg.ZeroOut(s.Log)
		}

	}

}
