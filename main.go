package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Project []struct {
		Name    string `yaml:"name"`
		Secrets []string
	} `yaml:"project"`
}

var config ConfigStruct
var log *logrus.Logger
var configpath string
var Version = "development"
var gitRef = "" // needs to be overridden in CI

func init() {
	// init logging
	log = logrus.New()
	log.Level = logrus.DebugLevel

	// set Version if gitRef is available
	gitRef = strings.TrimSpace(gitRef)
	refParts := strings.Split(gitRef, "/")
	if len(gitRef) > 0 && len(refParts) > 0 {
		Version = refParts[len(refParts)-1]
		log.Level = logrus.InfoLevel // also set log level to Info
	}

	// check if running on a supported OS
	useros := runtime.GOOS
	if useros == "darwin" {
		log.Debug("user OS is mac")
	} else {
		panic("not implemented")
	}

	// get home dir
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	home := user.HomeDir
	log.Debug("home: " + home)

	// get config folder
	configfolder := filepath.Join(home, ".config/uses")
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
		log.Fatal(err)
	}
	log.Tracef("Unmarshalled Config YAML:\n%v\n", config)
}

func main() {
	enablecli()
}
