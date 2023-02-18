package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	log "github.com/malikashish8/uses/logging"
	"gopkg.in/yaml.v3"
)

type ConfigSecret struct {
	Key          string `yaml:"key"`
	VariableName string `yaml:"variableName"`
}
type ConfigStruct struct {
	Project []struct {
		Name         string         `yaml:"name"`
		ConfigSecret []ConfigSecret `yaml:"secrets"`
	} `yaml:"project"`
}

var config ConfigStruct
var configpath string
var Version = "development"
var gitTag = "" // needs to be overridden in CI `go build -ldflags="main.gitTag=v0.1.1)"`

func init() {

	// set Version if gitTag is available
	if len(gitTag) > 0 {
		Version = gitTag
	}

	// check if running on a supported OS
	useros := runtime.GOOS
	if useros != "darwin" {
		panic("not implemented")
	}

	// get home dir
	user, err := user.Current()
	if err != nil {
		log.Fatal("Error: %s", err)
	}
	home := user.HomeDir

	// get config folder
	configfolder := filepath.Join(home, ".uses")
	configpath = filepath.Join(configfolder, "config.yaml")

	// init config if not exits
	_, err = os.Stat(configfolder)
	if err != nil {
		log.Info("Creating default config")
		os.Mkdir(configfolder, 0750)
	}
	_, err = os.Stat(configpath)
	if err != nil {
		_, err = os.Create(configpath)
		if err != nil {
			log.Fatal("Unable to create config file at " + configpath)
		}
	}

	// read config
	configyaml, err := ioutil.ReadFile(configpath)
	if err != nil {
		log.Fatal("unable to read config from " + configpath)
	}
	config = ConfigStruct{}
	err = yaml.Unmarshal(configyaml, &config)
	if err != nil {
		log.Fatal("Error: %s", err)
	}
	log.Trace("Unmarshalled Config YAML:\n%v\n", config)
}

func main() {
	enablecli()
}

// Custom unmarshaler for ConfigSecret to support `key as variableName` syntax
func (c *ConfigSecret) UnmarshalYAML(value *yaml.Node) error {
	var rawInput string
	// try to decode as a string and then check if it has ` as `
	if err := value.Decode(&rawInput); err == nil {
		parts := strings.Split(rawInput, " as ")
		if len(parts) == 2 {
			c.Key = parts[0]
			c.VariableName = parts[1]
			log.Debug("UnmarshalYAML successful with 2 values for: %v", c)
			return nil
		} else if len(parts) == 1 {
			c.Key = rawInput
			c.VariableName = rawInput
			log.Debug("UnmarshalYAML successful with 1 value for: %v", c)
			return nil
		} else {
			log.Error("UnmarshalYAML failed with value: %v", rawInput)
			return errors.New("invalid config for secrets at " + configpath + ": " + rawInput)
		}
	}

	// error handling - convert value to readable string and return error
	buffer := &bytes.Buffer{}
	for _, v := range value.Content {
		buffer.WriteString("\t" + v.Value)
	}

	return fmt.Errorf("invalid config for secrets at %v: %v", configpath, buffer)
}
