/*

useradd -r -M -s /bin/false analyze-goguy
mkdir -p /webproject/analyze-goguy
*/

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"
)

var main_file = "/webproject/analyze-goguy/var/mem.log"
var old_file = "/webproject/analyze-goguy/var/old/mem.log"
var old_dir = "/webproject/analyze-goguy/var/old"

/* Test
var main_file = "/home/mchirico/analyzeSystem/var/proc.log"
var old_file =  "/home/mchirico/analyzeSystem/var/old/proc.log"
var old_dir = "/home/mchirico/analyzeSystem/var/old"
*/

func zeroOut() {

	_ = os.Remove(main_file)
	f, err := os.OpenFile(main_file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	fmt.Println("Zeroed out file\n")
	fi, err := f.Stat()
	if err != nil {
		log.Panic(err)
	} else {
		fmt.Printf("file size:%v\n", fi.Size())

	}
}

func main() {

	os.Mkdir(old_dir, 0777)

	os.Rename(main_file, old_file)
	fmt.Println("Starting mem.go")

	for {
		size := logProcess()
		//fmt.Println(size)
		time.Sleep(3 * time.Second)
		if size > 1533791 {
			zeroOut()
		}

	}

}
