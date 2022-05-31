package main

import (
	"encoding/json"

	"github.com/99designs/keyring"
	"golang.org/x/exp/slices"
)

type Secret struct {
	Key   string
	Value string
}

var ring keyring.Keyring
var secrets []Secret

func init() {
	var err error
	ring, err = keyring.Open(keyring.Config{
		ServiceName: "uses",
	})
	if err != nil {
		log.Fatal(err)
	}
	err = readFromSecretStore()
	// if err is "item not found" then initialize secret in secret store
	if err == keyring.ErrKeyNotFound {
		log.Println("No secret found. Initializing...")
		log.Print(err)
		writeErr := writeToSecretStore()
		if writeErr != nil {
			log.Fatal(writeErr)
		}
	} else if err != nil {
		log.Fatal(err)
	}
}

func getSecretsString(name string) (value string, err error) {
	i, err := ring.Get(name)
	if err != nil {
		return "", err
	}
	return string(i.Data), nil
}

func saveSecretsString(name string, value string) (err error) {
	err = ring.Set(keyring.Item{
		Key:  name,
		Data: []byte(value),
	})
	return err
}

// read secrets from secret store
func readFromSecretStore() error {
	secretsString, err := getSecretsString("uses")
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(secretsString), &secrets)
	return err
}

// write secrets to secret store
func writeToSecretStore() error {
	secretsBytes, err := json.Marshal(secrets)
	if err != nil {
		return err
	}
	err = saveSecretsString("uses", string(secretsBytes))
	return err
}

// add a secret to secrets and save it
func AddSecret(newSecret Secret) error {
	secrets = append(secrets, newSecret)
	err := writeToSecretStore()
	return err
}

// get value of a secret
func GetSecretValue(key string) (string, error) {
	index := slices.IndexFunc(secrets, func(s Secret) bool { return s.Key == key })
	if index == -1 {
		return "", keyring.ErrKeyNotFound
	} else {
		value := secrets[index].Value
		return value, nil
	}
}

// remove secret
func DeleteSecret(key string) error {
	index := slices.IndexFunc(secrets, func(s Secret) bool { return s.Key == key })
	if index == -1 {
		return keyring.ErrKeyNotFound
	} else {
		secrets = slices.Delete(secrets, index, index)
		err := writeToSecretStore()
		return err
	}
}

// list secrets
func ListSecretKeys() []string {
	keys := make([]string, len(secrets))
	for i := 0; i < len(secrets); i++ {
		keys[i] = secrets[i].Key
	}
	return keys
}

// check if secrets has Key
func HasSecretKey(str string) bool {
	keys := ListSecretKeys()
	for _, v := range keys {
		if v == str {
			return true
		}
	}
	return false
}
