package yamlpkg

// REF: http://sweetohm.net/article/go-yaml-parsers.en.html

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
)

type y struct {
	Command         string `yaml:"Command"`
	Log             string `yaml:"Log"`
	LogSizeLimit    int    `yaml:"LogSizeLimit"`
	ArchiveLog      string `yaml:"ArchiveLog"`
	LoopDelay       int    `yaml:"LoopDelaySeconds"`
	DieAfterHours   int    `yaml:"DieAfterNumberHours"`
	DieAfterSeconds int    `yaml:"DieAfterNumberSeconds"`
}

// Config entry point
type Config struct {
	sync.Mutex
	Yaml y
}

// Write yaml file
func (c *Config) Write(file string) error {
	c.Lock()
	defer c.Unlock()

	data, err := yaml.Marshal(c.Yaml)
	if err != nil {
		log.Printf("yaml.Marshal(config): %v", err)
		return err
	}

	err = ioutil.WriteFile(file, data, 0600)
	if err != nil {
		log.Printf("error in yaml write: %v", err)
	}
	return err
}

// Read yaml file
func (c *Config) Read(file string) error {
	c.Lock()
	defer c.Unlock()

	source, err := ioutil.ReadFile(file)
	if err != nil {
		log.Printf("Error ioutil.ReadFile")
		return err
	}
	err = yaml.Unmarshal(source, &c.Yaml)
	if err != nil {
		log.Printf("Error Unmarshal")
		return err
	}

	return err
}

// SetDefault simple config settings
func (c *Config) SetDefault() {
	c.Lock()
	defer c.Unlock()
	c.Yaml.Command = `body() { IFS= read -r header; printf 'LogWrite: %s %s\n %s\n' $(date "+%Y-%m-%d %H:%M:%S") "$header"; "$@"; } && ps aux| body sort -n -r -k 4`
	c.Yaml.Log = "mem.log"
	c.Yaml.LoopDelay = 20
	c.Yaml.LogSizeLimit = 4000000
	c.Yaml.ArchiveLog = "memarchive.log"
	c.Yaml.DieAfterHours = 200
	c.Yaml.DieAfterSeconds = 2
}
