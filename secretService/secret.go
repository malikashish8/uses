package secretService

import (
	log "github.com/malikashish8/uses/logging"

	"github.com/99designs/keyring"
)

var ring keyring.Keyring

type Secret struct {
	Key   string // unique key for this secret
	Value string // value of this secret
}

func init() {
	if ring == nil {
		var err error
		// this always opens even if they keyring doesn't exist
		ring, err = keyring.Open(keyring.Config{
			ServiceName:                    "uses",
			KeychainName:                   "uses",
			KWalletAppID:                   "uses",
			KWalletFolder:                  "uses",
			KeychainTrustApplication:       true,
			WinCredPrefix:                  "uses",
			KeychainAccessibleWhenUnlocked: true,
		})
		if err != nil {
			log.Error("Unable to open keyring: %v", err)
		}
	}
}

// check if secret exists in keyring - does not require a password in MacOS keychain
func SecretExists(key string) (exists bool, err error) {
	keys, err := ring.Keys()
	if err != nil {
		return false, err
	}
	// print keys
	log.Debug("keys=%v\n", keys)

	// return error if key does not exit in keys
	for _, v := range keys {
		// print v and key

		log.Debug("v=%v, key=%v\n", v, key)
		if v == key {
			return true, nil
		}
	}
	return false, keyring.ErrKeyNotFound
}

func GetSecretValue(key string) (value string, err error) {
	// check if key exists in the keyring beforehand as secretExists() does not require a password
	exists, err := SecretExists(key)
	log.Debug("exists=%v, err=%v", exists, err)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", keyring.ErrKeyNotFound
	}

	item, err := ring.Get(key)
	if err != nil {
		return "", err
	}
	return string(item.Data), nil
}

// careful with this function, it will overwrite the existing secrets
func SaveSecret(key string, value string) (err error) {
	err = ring.Set(keyring.Item{
		Label: key,
		Key:   key,
		Data:  []byte(value),
	})
	return err
}

// remove secret
func DeleteSecret(key string) error {
	exists, err := SecretExists(key)
	if err != nil {
		return err
	}
	if !exists {
		return keyring.ErrKeyNotFound
	}
	err = ring.Remove(key)
	if err != nil {
		return err
	}
	return nil
}

// list secrets
func ListSecretKeys() ([]string, error) {
	return ring.Keys()
}
