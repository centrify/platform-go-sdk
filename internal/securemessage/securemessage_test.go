package securemessage

import (
	"flag"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

// secret test DOES NOT need to parse command line argument
// do this here to avoid errors if -config is passed in
// Anyway we need to declare them here so that go test will not complain
var (
	configPtr      = flag.String("config", "", "configuration file")
	configString   = flag.String("config-string", "", "configuration string")
	VaultRootToken = flag.String("vault-root-token", "root", "Vault root token")
)

type SecureMessageTestSuite struct {
	suite.Suite
}

func TestSecureMessageTestSuite(t *testing.T) {
	suite.Run(t, &SecureMessageTestSuite{})
}

func (s *SecureMessageTestSuite) TestSimpleMessages() {
	type testCase struct {
		msg   string
		label string
	}

	testCases := []testCase{
		{"test1", ""},        // no label
		{"test1", "label 1"}, // has label
		{"", "empty string"}, // empty string
	}

	key, _, err := GetPublicKey()
	s.Assert().NoError(err, "Should not have error in generating keys")

	for _, test := range testCases {
		res, err := EncryptString(test.msg, test.label, key)

		s.Assert().NoError(err, "Should not have error in encrypting message")
		s.T().Logf("Test message: <%s>, label: <%s> No. of blobs: %d\nBlob:<%v>\n", test.msg, test.label, len(res), res)

		// try decrypt the encrypted message
		clearText, err := DecryptString(res, test.label)
		s.Assert().NoError(err, "should have no error during decryption")
		s.Assert().Equal(test.msg, clearText, "Decrypted string should match")
	}
}

func (s *SecureMessageTestSuite) TestEncryptionTime() {
	txt := "This is a normal size data for encryption"
	label := "time test"

	key, _, err := GetPublicKey()
	s.Assert().NoError(err, "Should not have error in generating keys")

	t := time.Now()
	res, err := EncryptString(txt, label, key)
	s.Assert().NoError(err, "Should not have error in encrypting message")
	elapsed := time.Since(t)
	s.T().Logf("Elapsed time for encryption: %v microsecond\n", elapsed.Microseconds())

	t = time.Now()
	clearText, err := DecryptString(res, label)
	s.Assert().NoError(err, "should have no error during decryption")
	elapsed = time.Since(t)
	s.T().Logf("Elapsed time for decryption: %v microseconds\n", elapsed.Microseconds())
	s.Assert().Equal(txt, clearText, "should be able to decrypt")
}

func (s *SecureMessageTestSuite) TestMultipleChunks() {
	type testCase struct {
		msgLen    int
		expChunks int
	}

	testCases := []testCase{
		{msgLimit - 1, 1},   // one chunk
		{msgLimit, 1},       // exactly one chunk
		{msgLimit + 1, 2},   // two chunks
		{2*msgLimit - 1, 2}, // 2 chunks
		{2 * msgLimit, 2},   // 2 chunks
		{2*msgLimit + 1, 3}, // 3 chunks
	}

	key, _, err := GetPublicKey()
	s.Assert().NoError(err, "Should not have error in generating keys")

	var label string

	for _, test := range testCases {
		rawText := s.generateString(test.msgLen)
		label = fmt.Sprintf("Test message of length %d", len(rawText))
		res, err := EncryptString(rawText, label, key)

		s.Assert().NoError(err, "Should not have error in encrypting message")
		s.T().Logf("Test message: <%s>, label: <%s> No. of blobs: %d\nBlob:<%v>\n", rawText, label, len(res), res)
		s.Assert().Equal(test.expChunks, len(res), "number of chunks not expected")

		// try decrypt the encrypted message
		clearText, err := DecryptString(res, label)
		s.Assert().NoError(err, "should have no error during decryption")
		s.Assert().Equal(rawText, clearText, "Decrypted string should match")
	}
}
func (s *SecureMessageTestSuite) generateString(size int) string {

	var b strings.Builder
	var src = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789./~-=+*")
	srcLen := len(src)

	for i := 0; i < size; i++ {
		b.WriteRune(src[rand.Intn(srcLen)])
	}
	return b.String()
}
