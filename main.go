package main

import (
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/99designs/keyring"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Project []struct {
		Name    string `yaml:"name"`
		Secrets []string
	} `yaml:"project"`
}

var config ConfigStruct
var ring keyring.Keyring
var configpath string

func init() {
	// check if running on a supported os
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
	log.Debugf("Unmarshalled Config YAML:\n%v\n", config)

	// initalize keyring
	ring, err = keyring.Open(keyring.Config{
		ServiceName: "uses",
	})
	if err != nil {
		log.Fatal(err)
	}

}

func getSecret(name string) (value string, err error) {
	i, err := ring.Get(name)
	if err != nil {
		return "", err
	}
	return string(i.Data), nil
}

func setSecret(name string, value string) (err error) {
	err = ring.Set(keyring.Item{
		Key:  name,
		Data: []byte(value),
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	enablecli()
}
