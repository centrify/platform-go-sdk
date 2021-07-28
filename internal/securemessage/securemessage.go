// Package securemessage handles generation and use of one-shot asymmetric encrpyption keys used in LRPC messages
package securemessage

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	random "math/rand"
	"sync"
)

const keysize = 1024 // use 1024-bit RSA key for now... export limitation concern
const msgLimit = 62  // for 1024-bit RSA key using SHA-256

// common errors
var (
	ErrNoPublicKey = errors.New("No public key")
)

type keyInfo struct {
	mutex  sync.Mutex      // for synchronization
	hasKey bool            // whether key is generated or not
	err    error           // if not nil, key generation encountered error
	priv   *rsa.PrivateKey // private key
	keyID  uint32          // monotonic key ID.  Should update this if we want to rotate the private/public keypairs
}

var keyStore keyInfo

// GetPublicKey returns the public key for use in LRPC messages
func GetPublicKey() (*rsa.PublicKey, uint32, error) {
	err := checkKey()
	if err != nil {
		return nil, 0, err
	}
	return &keyStore.priv.PublicKey, keyStore.keyID, nil
}

// EncryptString encrpyts the payload using the given public key.  This is usually done by the LRPC client when it needs to
// send sensitive information (e.g., password) in LRPC message.  A slice of base64-encoded strings
// is returned as this function need to chunk the input string into one or more blocks for encryption.
func EncryptString(payload string, label string, pubKey *rsa.PublicKey) ([]string, error) {

	if pubKey == nil {
		return nil, ErrNoPublicKey
	}
	msglen := len(payload)
	chunks := (msglen-1)/msgLimit + 1

	if chunks < 1 {
		chunks = 1 // worst case of empty string
	}
	ret := make([]string, chunks)
	src := []byte(payload)
	index := 0

	// loop to work on a chunk at a time
	for chunk := 0; chunk < chunks; chunk++ {
		endIndex := index + msgLimit
		if endIndex > msglen {
			endIndex = msglen
		}
		cipher, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, pubKey, src[index:endIndex], []byte(label))
		if err != nil {
			return nil, err
		}
		ret[chunk] = base64.StdEncoding.EncodeToString(cipher)
		index = endIndex
	}

	return ret, nil
}

// DecryptString decrypts the slice of based64-encrypted strings into a single text string
func DecryptString(ciphers []string, label string) (string, error) {
	err := checkKey()
	if err != nil {
		return "", err
	}

	ret := ""
	for _, cipherText := range ciphers {
		cipher, err := base64.StdEncoding.DecodeString(cipherText)
		if err != nil {
			return "", err // not a good hex string
		}
		plain, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, keyStore.priv, cipher, []byte(label))
		if err != nil {
			return "", err // error in decrypt
		}
		ret = ret + string(plain)
	}

	return ret, nil
}

// checkKey checks if a key is generated or not, or whether it encounters error in key generation
func checkKey() error {
	keyStore.mutex.Lock()
	defer keyStore.mutex.Unlock()

	if keyStore.err != nil {
		return keyStore.err // return key generation error
	}
	if !keyStore.hasKey {
		// no key pair yet, generates it
		priv, err := rsa.GenerateKey(rand.Reader, keysize)
		if err != nil {
			keyStore.err = err
			return err
		}
		// private key generated
		keyStore.priv = priv
		keyStore.priv.Precompute() // for faster private key operations
		keyStore.hasKey = true
		keyStore.keyID = uint32(random.Int31())
	}
	return nil // success
}
