package main

import (
	"encoding/json"
	"sync"

	"github.com/99designs/keyring"
	"golang.org/x/exp/slices"
)

type SecretService struct {
	secretOnce sync.Once
	secrets    []Secret
}

type Secret struct {
	Key   string
	Value string
}

var secretService SecretService
var ring keyring.Keyring

// lazy initialize secret service so that secret is not access unless needed
func (s *SecretService) lazyInit() {
	s.secretOnce.Do(func() {
		var err error
		ring, err = keyring.Open(keyring.Config{
			ServiceName: "uses",
		})
		if err != nil {
			log.Fatal(err)
		}
		err = s.readFromSecretStore()
		// if err is "item not found" then initialize secret in secret store
		if err == keyring.ErrKeyNotFound {
			log.Println("No secret found. Initializing...")
			writeErr := s.writeToSecretStore()
			if writeErr != nil {
				log.Fatal(writeErr)
			}
		} else if err != nil {
			log.Fatal(err)
		}
	})
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
func (s *SecretService) readFromSecretStore() error {
	secretsString, err := getSecretsString("uses")
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(secretsString), &s.secrets)
	return err
}

// write secrets to secret store
func (s *SecretService) writeToSecretStore() error {
	secretsBytes, err := json.Marshal(s.secrets)
	if err != nil {
		return err
	}
	err = saveSecretsString("uses", string(secretsBytes))
	return err
}

// add a secret to secrets and save it
func (s *SecretService) AddSecret(newSecret Secret) error {
	s.lazyInit()
	s.secrets = append(s.secrets, newSecret)
	err := s.writeToSecretStore()
	return err
}

// get value of a secret
func (s *SecretService) GetSecretValue(key string) (string, error) {
	s.lazyInit()
	index := slices.IndexFunc(s.secrets, func(s Secret) bool { return s.Key == key })
	if index == -1 {
		return "", keyring.ErrKeyNotFound
	} else {
		value := s.secrets[index].Value
		return value, nil
	}
}

// remove secret
func (s *SecretService) DeleteSecret(key string) error {
	s.lazyInit()
	index := slices.IndexFunc(s.secrets, func(s Secret) bool { return s.Key == key })
	if index == -1 {
		return keyring.ErrKeyNotFound
	} else {
		s.secrets = slices.Delete(s.secrets, index, index+1)
		err := s.writeToSecretStore()
		return err
	}
}

// list secrets
func (s *SecretService) ListSecretKeys() []string {
	s.lazyInit()
	keys := make([]string, len(s.secrets))
	for i := 0; i < len(s.secrets); i++ {
		keys[i] = s.secrets[i].Key
	}
	return keys
}

// check if secrets has Key
func (s *SecretService) HasSecretKey(str string) bool {
	s.lazyInit()
	keys := s.ListSecretKeys()
	for _, v := range keys {
		if v == str {
			return true
		}
	}
	return false
}
