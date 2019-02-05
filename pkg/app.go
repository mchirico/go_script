package pkg

import (
	"context"
	"fmt"
	"github.com/mchirico/go_script/analyze"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Script is may way to these routines
type Script struct {
	sync.Mutex
	JSON    JSON
	Analyze analyze.A
}

// JSON read in from config
type JSON struct {
	Command       string `json:"command"`
	Log           string `json:"log"`
	LogSizeLimit  int    `json:"logSizeLimit"`
	ArchiveLog    string `json:"logArchive"`
	LoopDelay     int    `json:"loopDelay"`
	DieAfterHours int    `json:"dieAfterNumHours"`
}

// ReadConfig reads configuration
func (s *Script) ReadConfig() {

	s.JSON.Command = "body() { IFS= read -r header; printf '%s %s\n %s\n' `date \"+%Y-%m %H:%M:%S\"` \"$header\"; \"$@\"; } && ps aux| body sort -n -r -k 4 && free"
	s.JSON.Log = "/tmp/s.log"
	s.JSON.ArchiveLog = "/tmp/archive/s.log"

}

// Writer for future implementation
type Writer interface {
	Write(file string, data []byte) (n int, err error)
}

// MakeDir only makes single directory
func MakeDir(dir string) {
	os.Mkdir(dir, 0777)
	//os.Rename(main_file, old_file)
}

// ZeroOut file
func ZeroOut(file string) (int64, error) {

	_ = os.Remove(file)
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()

	if fi, err := f.Stat(); err != nil {
		return -1, err
	} else {
		f.WriteString(fmt.Sprintf("file size:%v\n", fi.Size()))
		return fi.Size(), err
	}
}

// ReadFile is a generic read function
func ReadFile(file string) (string, error) {
	data, err := ioutil.ReadFile(file)
	return string(data), err
}

// WriteData generic write function
func WriteData(file string, data []byte) (int, int64, error) {

	f, err := os.OpenFile(file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return -1, -1, err
	}
	defer f.Close()

	if fi, err := f.Stat(); err != nil {
		return 0, -1, err
	} else {
		n, err := f.WriteString(fmt.Sprintf("file size:%v\n%s", fi.Size(), data))
		if err == nil {
			log.Printf("wrote:%v\nfileSize:%v\n", n, fi.Size()+int64(n))
		}
		return n, fi.Size(), err
	}
}

// LogProcess return bytes written and total bytes in file
func (s *Script) LogProcess(ctx context.Context) (int, int64, []byte) {

	s.Lock()
	defer s.Unlock()
	//cmd := exec.Command("sh", "-c", "date '+%Y-%m-%d %H:%M:%S\n' 1>&2;top -b -n1 -c -w 400 -o +%MEM|head -n30 1>&2")
	cmd := exec.CommandContext(ctx, "sh", "-c", s.JSON.Command)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("LogProcess command: %s", s.JSON.Command)
		log.Fatalf("(fatal) LogProcess: %v", err)
	}

	n, fileSize, err := WriteData(s.JSON.Log, output)
	if err != nil {
		panic(err)
	}
	return n, fileSize, output

}

// Process - one call on loop
func (s *Script) Process(milliseconds time.Duration, sizeLimit int64) []byte {

	ctx, cancel := context.WithTimeout(context.Background(),
		milliseconds*time.Millisecond)
	defer cancel()

	_, size, output := s.LogProcess(ctx)
	if size > sizeLimit {
		ZeroOut(s.JSON.Log)
	}
	return output

}

func delay(delay int) {
	if delay <= 0 {
		time.Sleep(1 * time.Second)
	} else {
		time.Sleep(time.Duration(delay) * time.Second)
	}

}

// Loop through commands
func (s *Script) Loop(ctx context.Context, milliseconds time.Duration, sizeLimit int64) {

	gen := func(ctx context.Context) <-chan []byte {
		dst := make(chan []byte)
		var output []byte
		go func() {
			for {
				select {

				case <-ctx.Done():
					log.Printf("ctx.Done() in Loop")
					dst <- []byte{}
					return // returning not to leak the goroutine
				case dst <- output:
					output = s.Process(milliseconds, sizeLimit)
				}

				delay(s.JSON.LoopDelay)

			}
		}()
		return dst
	}

	for output := range gen(ctx) {

		if s.Analyze != nil {
			s.Analyze(output)
		}

		if ctx.Err() != nil {
			log.Printf("calling break")
			break
		}

	}

}
