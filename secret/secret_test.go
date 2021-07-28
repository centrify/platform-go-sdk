package secret

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/centrify/cloud-golang-sdk/oauth"
	"github.com/centrify/platform-go-sdk/testutils"
	"github.com/stretchr/testify/suite"
)

type SecretTestSuite struct {
	testutils.CfyTestSuite
	Token        string // access token
	handle       Secret // interface to secret API
	suffix       string // random suffix used in this test
	testFolder   string // path to test folder
	testFolderID string // test folder ID
}

/* Error conditions to test:
- create empty secret
- secret does not exist for operations
- secret already exist for create
- verify deleted folder/secret does not exist in get
*/

func TestSecretTestSuite(t *testing.T) {
	suite.Run(t, new(SecretTestSuite))
}

func (s *SecretTestSuite) SetupSuite() {
	var err error
	s.LoadConfig()

	if s.Config.TenantURL == "" {
		s.T().Skip("Tenant URL must be specified")
	}
	if !isPASTenantConnected(s.Config.TenantURL) {
		s.T().Skipf("Tenant %s is not online\n", s.Config.TenantURL)
	}

	// TODO: handle case of using DMC for testing
	if s.Config.PASuser.Username == "" || s.Config.PASuser.Password == "" {
		s.T().Skip("Must specify username and password for test PAS users")
	}
	if s.Config.AppID == "" || s.Config.Scope == "" {
		s.T().Skip("Must specify an web application and scope for test")
	}

	// get access token for the user
	s.Token, err = s.getAccessToken()
	s.Require().NoError(err, "Require Oauth token to continue testing")

	// get handle to secret interface for testing
	s.handle, err = NewSecretClient(s.Config.TenantURL, ServerPAS, s.Token, nil)
	s.Require().NoError(err, "Requires client handle for access to secrets")

	// generate random string as suffix for use in tests
	t := time.Now().UnixNano()
	s.suffix = fmt.Sprintf("-%x", t)
	s.T().Logf("Suffix used in test: [%s]\n", s.suffix)

	folderPath := "Test-folder" + s.suffix
	s.testFolderID = s.createTestFolder(folderPath)
	s.testFolder = folderPath
}

// delete the test folder on cleanup
func (s *SecretTestSuite) TearDownSuite() {
	if s.testFolder != "" {
		if s.handle != nil {
			s.handle.Delete(s.testFolder)
		} else {
			s.T().Error("No handle to clean up test folder.  Logical error in test")
		}
	}
}
func (s *SecretTestSuite) TestSimpleTextSecret() {
	type textTestCase struct {
		path    string
		content string
	}

	testcases := []textTestCase{
		{path: "top_folder_secret" + s.suffix,
			content: "test string with space",
		},
		{path: s.testFolder + "/textsecret1" + s.suffix,
			content: "test data with / and stored in secret in folder",
		},
		{
			path:    s.testFolder + "/textsecret2" + s.suffix,
			content: "yet another secret",
		},
		{
			path:    s.testFolder + "/emptystring" + s.suffix,
			content: "  ",
		},
	}

	var id string

	for _, testcase := range testcases {
		id = s.createTextSecret(testcase.path, testcase.content)
		// for file cleanup
		defer s.handle.Delete(testcase.path)

		readback, r, err := s.handle.Get(testcase.path)
		s.Assert().NoError(err, "Should not return error on readback")
		s.Assert().Equal(testcase.content, readback, "Read back result must be the same")
		s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

		// get metadata
		metadata, r, err := s.handle.GetMetaData(testcase.path)
		s.Assert().NoError(err, "Should not return error on getting metadata")
		s.Assert().Equal(id, metadata.ID, "Object ID must be the same")
		s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
	}
}

func (s *SecretTestSuite) TestSimpleKVSecret() {

	type kvTestCase struct {
		path    string
		content map[string]string
	}

	testcases := []kvTestCase{
		{
			path: "top_folder_kvsecret" + s.suffix,
			content: map[string]string{
				"location": "top level folder",
				"foo":      "bar",
			},
		},
		{
			path: s.testFolder + "/testkv1" + s.suffix,
			content: map[string]string{
				"foo":   "bar",
				"hello": "world",
			},
		},
		{
			path:    s.testFolder + "/emptykv" + s.suffix,
			content: map[string]string{},
		},
		{
			path: s.testFolder + "empty_value" + s.suffix,
			content: map[string]string{
				"foo": "",
			},
		},
		{
			path: s.testFolder + "empty_key" + s.suffix,
			content: map[string]string{
				"": "no key",
			},
		},
	}

	var id string

	for _, testcase := range testcases {
		id = s.createKVSecret(testcase.path, testcase.content)
		// for file cleanup
		defer s.handle.Delete(testcase.path)

		readback, r, err := s.handle.Get(testcase.path)
		s.Assert().NoError(err, "Should not return error on readback")
		s.Assert().Equal(testcase.content, readback, "Read back result must be the same")
		s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

		// get metadata
		metadata, r, err := s.handle.GetMetaData(testcase.path)
		s.Assert().NoError(err, "Should not return error on getting metadata")
		s.Assert().Equal(id, metadata.ID, "Object ID must be the same")
		s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
	}
}

func (s *SecretTestSuite) TestListEmptyFolder() {

	path := "emptyFolder" + s.suffix
	id := s.createTestFolder(path)
	defer s.handle.Delete(path)

	// now list the folder
	items, r, err := s.handle.List(path)
	s.Assert().NoError(err, "Should not return error in getting content of empty folder")
	s.Assert().NotEmpty(id, "must return object ID")
	s.Assert().Len(items, 0, "Should have no entries returned")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

	// get metadata
	metadata, r, err := s.handle.GetMetaData(path)
	s.Assert().NoError(err, "Should not return error on getting metadata")
	s.Assert().Equal(id, metadata.ID, "Object ID must be the same")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
}

func (s *SecretTestSuite) TestListTopLevelFolder() {

	// test will create a text secret
	// since suite setup also create a test folder, we should have at
	// least 2 items
	secretPath := "listing_test" + s.suffix
	id := s.createTextSecret(secretPath, "something")
	defer s.handle.Delete(secretPath)

	// setup the expected items
	expItems := []Item{
		{Name: secretPath, Type: SecretTypeText, ID: id},
		{Name: s.testFolder, Type: SecretTypeFolder, ID: s.testFolderID},
	}

	items, r, err := s.handle.List("/")
	s.Assert().NoError(err, "Should not get error when listing top level folder")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
	s.Assert().GreaterOrEqual(len(items), len(expItems), "Should return at least what we created")

	for _, expItem := range expItems {
		s.Assert().Contains(items, expItem)
	}
}

// try to create a secret that does not exist in top level folder yet
// Will retry till one is found...
func (s *SecretTestSuite) findNonExistentSecretInTop(path string) (string, error) {

	trial := 1
	fpath := path

	for {
		_, _, err := s.handle.GetMetaData(fpath)
		if err == ErrSecretNotFound {
			// this is the one
			return fpath, nil
		}
		fpath = fmt.Sprintf("%s_%d", path, trial)
		trial++
		if trial > 100 {
			s.T().Log("Too many path conflicts in top level folder")
			return "", errors.New("Too many path conflicts")
		}
	}
}
func (s *SecretTestSuite) TestTopLevelSecret() {

	content := "This is test content"

	path := "toptest" + s.suffix
	fpath, err := s.findNonExistentSecretInTop(path)
	if err != nil {
		s.T().Skipf("Top level folder has too many secrets that conflict with test prefix [%s]", path)
	}

	s.createTextSecret(fpath, content)
	defer func() {
		_, err := s.handle.Delete(fpath)
		s.Assert().NoError(err, "Should have no error in deleting the secret")
	}()

	readback, r, err := s.handle.Get(fpath)
	s.Assert().NoError(err, "should not return error on readback")
	s.Assert().Equal(content, readback, "Read back result should be the same")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

	// verify that the item is in result when folder is listed
	items, r, err := s.handle.List("/")
	s.Assert().NoError(err, "Should not get error when listing top level folder")
	s.Assert().Equal(200, r.StatusCode, "HTTP Status should be 200")
	s.Assert().GreaterOrEqual(len(items), 1, "Should return at least one item")
	// result should contain the object
	found := false
	for _, item := range items {
		if item.Name == fpath && item.Type == SecretTypeText {
			found = true
			break
		}
	}
	s.Assert().True(found, "list result must contain created item")

}
func (s *SecretTestSuite) TestListFolder() {
	t := s.T()

	// test will create a separate folder for testing
	// so that we can have the exact count

	secretFolder := "listing_content_test" + s.suffix
	id := s.createTestFolder(secretFolder)
	defer s.handle.Delete(secretFolder)

	// test cases
	type testcase struct {
		path        string
		secretType  string
		textContent string
		kvContent   map[string]string
	}
	// setup test cases
	testcases := []testcase{
		{
			path:       "kv1",
			secretType: SecretTypeKV,
			kvContent: map[string]string{
				"location": "top level folder",
				"foo":      "bar",
			},
		},
		{
			path:       "kv2",
			secretType: SecretTypeKV,
			kvContent: map[string]string{
				"foo":   "bar",
				"hello": "world",
			},
		},
		{
			path:       "kv3",
			secretType: SecretTypeKV,
			kvContent:  map[string]string{},
		},
		{
			path:        "secret1",
			secretType:  SecretTypeText,
			textContent: "test string with space",
		},
		{
			path:        "secret2",
			secretType:  SecretTypeText,
			textContent: "test data with / and stored in secret in folder",
		},
		{
			path:        "secret3",
			secretType:  SecretTypeText,
			textContent: "yet another secret",
		},
		{
			path:       "subfolder1",
			secretType: SecretTypeFolder,
		},
		{
			path:       "subfolder2",
			secretType: SecretTypeFolder,
		},
	}

	expItems := make([]Item, len(testcases))
	var path string

	// create the items in the test cases
	for i, tc := range testcases {
		path = secretFolder + "/" + tc.path
		switch tc.secretType {
		case SecretTypeFolder:
			id = s.createTestFolder(path)
		case SecretTypeText:
			id = s.createTextSecret(path, tc.textContent)
		case SecretTypeKV:
			id = s.createKVSecret(path, tc.kvContent)
		default:
			t.Errorf("*** error in test case setup. Unknow secretType: %s\n", tc.secretType)
		}
		defer s.handle.Delete(path)
		// setup the expected values in the directory listing
		expItems[i] = Item{
			Name: tc.path,
			Type: tc.secretType,
			ID:   id,
		}
	}

	items, r, err := s.handle.List(secretFolder)
	s.Assert().NoError(err, "Should not get error when listing top level folder")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
	s.Assert().Equal(len(expItems), len(items))
	for _, expItem := range expItems {
		s.Assert().Contains(items, expItem)
	}

}

func (s *SecretTestSuite) TestModifyTextSecret() {

	path := s.testFolder + "/modify_secret_test"
	origContent := "this is original"
	newContent := "this is new content"

	id := s.createTextSecret(path, origContent)
	// for file cleanup
	defer s.handle.Delete(path)

	s.testModifySingleTextSecret(path, origContent, newContent, id)
}

func (s *SecretTestSuite) TestModifyKVSecret() {

	path := s.testFolder + "/modify_kv_test"
	origContent := map[string]string{
		"foo":   "bar",
		"hello": "world",
	}
	newContent := map[string]string{
		"bar":   "foo",
		"world": "hello",
	}

	id := s.createKVSecret(path, origContent)
	// for file cleanup
	defer s.handle.Delete(path)

	// get metadata
	origMetadata, r, err := s.handle.GetMetaData(path)
	s.Assert().NoError(err, "Should not return error on getting metadata")
	s.Assert().Equal(id, origMetadata.ID, "Object ID must be the same")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

	// now modify the data
	success, id2, r, err := s.handle.Modify(path, "", newContent)
	s.Assert().NoError(err, "Should not return error on getting metadata")
	s.Assert().True(success, "Modification should be successful")
	s.Assert().Equal(id, id2, "Object ID must be the same")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

	readback, r, err := s.handle.Get(path)
	s.Assert().NoError(err, "Should not return error on readback")
	s.Assert().Equal(newContent, readback, "Read back result must be the modified data")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

	// get metadata
	metadata, r, err := s.handle.GetMetaData(path)
	s.Assert().NoError(err, "Should not return error on getting metadata")
	s.Assert().Equal(id, metadata.ID, "Object ID must be the same")
	s.Assert().Equal(origMetadata.WhenCreated, metadata.WhenCreated, "creation time should not be changed")
	s.Assert().Equal(origMetadata.CRN, metadata.CRN, "CRN should not be changed")
	s.Assert().NotNil(metadata.WhenModified, "Should now have modification time")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

}

// Negative tests
func (s *SecretTestSuite) TestCreateSecretWhenExist() {

	origContent := "test original content"
	path := s.testFolder + "/preexist_secret"
	id := s.createTextSecret(path, origContent)
	defer s.handle.Delete(path)

	// try to create again
	success, _, r, err := s.handle.Create(path, "", "new content")
	s.Assert().ErrorIs(err, ErrExists, "Expects ErrExists")
	s.Assert().False(success, "Create should return false")
	s.Assert().Equal(409, r.StatusCode, "HTTP status should be 409")

	// try to create folder
	success, _, r, err = s.handle.CreateFolder(path, "")
	s.Assert().ErrorIs(err, ErrExists, "Expects ErrExists")
	s.Assert().False(success, "Create should return false")
	s.Assert().Equal(409, r.StatusCode, "HTTP status should be 409")

	// getMetadata should return the original ID
	metadata, r, err := s.handle.GetMetaData(path)
	s.Assert().NoError(err, "Should not return error when getting metadata")
	s.Assert().Equal(id, metadata.ID, "should have same ID")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

	// get should return the original data
	readback, r, err := s.handle.Get(path)
	s.Assert().NoError(err, "Should not return error on get")
	s.Assert().Equal(origContent, readback, "Should read back original content")
	s.Assert().Equal(200, r.StatusCode, "HTTP status code should be 200")

}

func (s *SecretTestSuite) TestCreateSecretWhenExistAsFolder() {

	path := s.testFolder + "/preexist_secret_folder"
	id := s.createTestFolder(path)
	defer s.handle.Delete(path)

	// try to create again as secret
	success, _, r, err := s.handle.Create(path, "", "new content")
	s.Assert().ErrorIs(err, ErrExists, "Expects ErrExists")
	s.Assert().False(success, "Create should return false")
	s.Assert().Equal(409, r.StatusCode, "HTTP status should be 409")

	// try to create folder
	success, _, r, err = s.handle.CreateFolder(path, "")
	s.Assert().ErrorIs(err, ErrExists, "Expects ErrExists")
	s.Assert().False(success, "Create should return false")
	s.Assert().Equal(409, r.StatusCode, "HTTP status should be 409")

	// getMetadata should return the original ID
	metadata, r, err := s.handle.GetMetaData(path)
	s.Assert().NoError(err, "Should not return error when getting metadata")
	s.Assert().Equal(id, metadata.ID, "should have same ID")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
}
func (s *SecretTestSuite) TestCreateSecretWithExistingFolderInPath() {

	path := s.testFolder + "/preexist_secret_folder"
	id := s.createTestFolder(path)
	defer s.handle.Delete(path)

	// try to create secret in same path
	success, _, r, err := s.handle.Create(path, "", "new content")
	s.Assert().ErrorIs(err, ErrExists, "Expects ErrExists")
	s.Assert().False(success, "Create should return false")
	s.Assert().Equal(409, r.StatusCode, "HTTP status should be 409")

	// getMetadata should return the original ID
	metadata, r, err := s.handle.GetMetaData(path)
	s.Assert().NoError(err, "Should not return error when getting metadata")
	s.Assert().Equal(id, metadata.ID, "should have same ID")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
}

// TestCreateSecretWithNoParent creates a secret where the parent folder does not
// exist
func (s *SecretTestSuite) TestCreateTextSecretWithNoParent() {

	folderPath := s.testFolder + "/parent_does_not_exist"
	obj := "secret_text"
	path := folderPath + "/" + obj
	content := "data in secret"

	// try to create secret directly
	s.createTextSecret(path, content)
	defer func() {
		// cleanup...
		s.handle.Delete(path)
		_, err := s.handle.Delete(folderPath)
		s.Assert().NoError(err, "Should not have error in deleting parent folder")
	}()

	// make sure that we can list the content of parent
	items, r, err := s.handle.List(folderPath)
	s.Assert().NoError(err, "Should not get error when listing parent folder")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
	s.Assert().Equal(1, len(items)) // should return just one item

	s.Assert().Equal(obj, items[0].Name, "should return name of created secret")
	s.Assert().Equal(SecretTypeText, items[0].Type, "type should be text")

	// readback
	readback, r, err := s.handle.Get(path)
	s.Assert().NoError(err, "should not return error on readback")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
	s.Assert().Equal(content, readback, "content should be the same")
}

func (s *SecretTestSuite) TestCreateKVSecretWithNoParent() {

	folderPath := s.testFolder + "/parent2_does_not_exist"
	obj := "secret_kv"
	path := folderPath + "/" + obj
	content := map[string]string{
		"bar":      "foo",
		"location": "no parent",
	}

	// try to create secret directly
	s.createKVSecret(path, content)
	defer func() {
		// cleanup...
		s.handle.Delete(path)
		_, err := s.handle.Delete(folderPath)
		s.Assert().NoError(err, "Should not have error in deleting parent folder")
	}()

	// make sure that we can list the content of parent
	items, r, err := s.handle.List(folderPath)
	s.Assert().NoError(err, "Should not get error when listing parent folder")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
	s.Assert().Equal(1, len(items)) // should return just one item

	s.Assert().Equal(obj, items[0].Name, "should return name of created secret")
	s.Assert().Equal(SecretTypeKV, items[0].Type, "type should be keyvalue")

	// readback
	readback, r, err := s.handle.Get(path)
	s.Assert().NoError(err, "should not return error on readback")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
	s.Assert().Equal(content, readback, "content should be the same")
}

// negative tests

func (s *SecretTestSuite) TestNonExistentSecret() {
	path := s.testFolder + "/does_not_exist"

	_, _, err := s.handle.GetMetaData(path)
	s.Assert().ErrorIs(err, ErrSecretNotFound, "Should get ErrSecretNotFound when secret does not exist")

	_, _, getErr := s.handle.Get(path)
	s.Assert().ErrorIs(getErr, ErrSecretNotFound, "Should get ErrSecretNotFound when getting secret")

	success, _, _, modErr := s.handle.Modify(path, "", "new content")
	s.Assert().False(success, "Modification should fail for non-existent secret")
	s.Assert().ErrorIs(modErr, ErrSecretNotFound, "Should get ErrSecretNotFound on modify")

	_, delErr := s.handle.Delete(path)
	s.Assert().ErrorIs(delErr, ErrSecretNotFound, "Should get ErrSecretNotFound on delete")

}
func (s *SecretTestSuite) TestModifySecretFolder() {

	path := s.testFolder + "/test_modify_folder"
	_ = s.createTestFolder(path)
	defer s.handle.Delete(path)

	// try to modify it
	success, _, _, err := s.handle.Modify(path, "", "new content")
	s.Assert().False(success, "Modification of folder should fail")
	s.Assert().Error(err, "Error should be returned")
	s.Assert().ErrorIs(err, ErrCannotModifySecretFolder, "Err should be CannotModifySecretFolder")

}

func (s *SecretTestSuite) TestModifyNonExistentSecret() {

	path := s.testFolder + "/non-existent-secret"

	// verify that it does not exist
	_, r, err := s.handle.GetMetaData(path)
	s.Assert().ErrorIs(err, ErrSecretNotFound, "Should get error when trying to retrieve metadata of non-existing secret")
	s.Assert().Equal(404, r.StatusCode, "HTTP status should be 404")

	// try to modify it
	success, _, r, err := s.handle.Modify(path, "", "new content")
	s.Assert().False(success, "Modification should fail")
	s.Assert().ErrorIs(err, ErrSecretNotFound)
	s.Assert().Equal(404, r.StatusCode, "HTTP status should be 404")
}
func (s *SecretTestSuite) TestChangeTextSecretToKV() {

	origContent := "test original content"
	path := s.testFolder + "/modify_text_secret"
	_ = s.createTextSecret(path, origContent)
	defer s.handle.Delete(path)

	// try to modify it again
	newContent := map[string]string{
		"foo":   "bar",
		"hello": "world",
	}
	success, _, r, err := s.handle.Modify(path, "", newContent)
	s.Assert().False(success, "Modification should fail")
	s.Assert().Error(err, "Should return error")
	s.Assert().ErrorIs(err, ErrCannotModifySecretType, "Should get CannotModifySecretType error")

	// get should return the original data
	readback, r, err := s.handle.Get(path)
	s.Assert().NoError(err, "Should not return error on get")
	s.Assert().Equal(origContent, readback, "Should read back original content")
	s.Assert().Equal(200, r.StatusCode, "HTTP status code should be 200")

}

func (s *SecretTestSuite) TestDeleteNonExistentSecret() {
	path := s.testFolder + "/non-existent-secret-for-delete"

	// verify that it does not exist
	_, r, err := s.handle.GetMetaData(path)
	s.Assert().ErrorIs(err, ErrSecretNotFound, "Should get error when trying to retrieve metadata of non-existing secret")
	s.Assert().Equal(404, r.StatusCode, "HTTP status should be 404")

	// try to delete it
	r, err = s.handle.Delete(path)
	s.Assert().ErrorIs(err, ErrSecretNotFound, "Should return ErrSecretNotFound")
	s.Assert().Equal(404, r.StatusCode, "HTTP status should be 404")
}

func (s *SecretTestSuite) TestDeleteFolderWhenNotEmpty() {
	folder := s.testFolder + "/non_empty_folder_to_delete"

	// create test folder first
	s.createTestFolder(folder)
	defer s.handle.Delete(folder) // don't care whether this fails or not

	// create a test secret inside folder
	secretPath := folder + "/a_secret"
	s.createTextSecret(secretPath, "just some content")
	defer s.handle.Delete(secretPath)

	// now try to delete the test folder
	r, err := s.handle.Delete(folder)
	s.Assert().ErrorIs(err, ErrFolderNotEmpty)
	s.Assert().Equal(409, r.StatusCode, "HTTP status code should be 409")

}
func (s *SecretTestSuite) TestChangeKVSecretToText() {
	origContent := map[string]string{
		"foo":   "bar",
		"hello": "world",
	}
	path := s.testFolder + "/modify_kv_secret"
	_ = s.createKVSecret(path, origContent)
	defer s.handle.Delete(path)

	// try to modify it again
	newContent := "new text content"
	success, _, r, err := s.handle.Modify(path, "", newContent)
	s.Assert().False(success, "Modification should fail")
	s.Assert().Error(err, "Should return error")
	s.Assert().ErrorIs(err, ErrCannotModifySecretType, "Should get CannotModifySecretType error")

	// get should return the original data
	readback, r, err := s.handle.Get(path)
	s.Assert().NoError(err, "Should not return error on get")
	s.Assert().Equal(origContent, readback, "Should read back original content")
	s.Assert().Equal(200, r.StatusCode, "HTTP status code should be 200")

}

func (s *SecretTestSuite) testModifySingleTextSecret(path, originalContent, newContent, id string) {
	// get metadata
	origMetadata, r, err := s.handle.GetMetaData(path)
	s.Assert().NoError(err, "Should not return error on getting metadata")
	s.Assert().Equal(id, origMetadata.ID, "Object ID must be the same")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

	// now modify the data
	success, id2, r, err := s.handle.Modify(path, "", newContent)
	s.Assert().NoError(err, "Should not return error on getting metadata")
	s.Assert().True(success, "Modification should be successful")
	s.Assert().Equal(id, id2, "Object ID must be the same")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

	readback, r, err := s.handle.Get(path)
	s.Assert().NoError(err, "Should not return error on readback")
	s.Assert().Equal(newContent, readback, "Read back result must be the modified data")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

	// get metadata
	metadata, r, err := s.handle.GetMetaData(path)
	s.Assert().NoError(err, "Should not return error on getting metadata")
	s.Assert().Equal(id, metadata.ID, "Object ID must be the same")
	s.Assert().Equal(origMetadata.WhenCreated, metadata.WhenCreated, "creation time should not be changed")
	s.Assert().Equal(origMetadata.CRN, metadata.CRN, "CRN should not be changed")
	s.Assert().NotNil(metadata.WhenModified, "Should now have modification time")
	s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")

}
func (s *SecretTestSuite) TestCreateInvalidSecretName() {
	type nameTestCase struct {
		name     string
		expOk    bool
		expError error
	}

	testContent := "this is test content"
	newContent := "this is modified"
	testcases := []nameTestCase{
		{name: "---", expOk: true, expError: nil},
		{name: "&#", expOk: false, expError: ErrBadPathName},
		{name: "abc<def", expOk: false, expError: ErrBadPathName},
		{name: "", expOk: false, expError: ErrExists},
		{name: "   ", expOk: false, expError: ErrBadPathName},

		{name: "  /empty_parent", expOk: false, expError: ErrBadPathName},

		// The following test cases are skipped now due to DV-231

		/*
			{name: "<", expOk: true, expError: nil},
			{name: ">", expOk: true, expError: nil},
			{name: "def>abc", expOk: true, expError: nil},
			{name: "&", expOk: true, expError: nil},
		*/
	}

	var secretName string
	for _, tc := range testcases {
		secretName = s.testFolder + "/" + tc.name
		s.T().Logf("For path [%s]:\n", secretName)
		success, id, _, err := s.handle.Create(secretName, "", testContent)
		if success {
			defer s.handle.Delete(secretName)
		}
		if tc.expOk {
			// expects success
			s.Assert().True(success, "Expects success")
			s.Assert().NoError(err, "Expects no error")

			// verify that I can read/modify/read/getMetadata the secret
			s.testModifySingleTextSecret(secretName, testContent, newContent, id)

		} else {
			// expects error
			s.Assert().False(success, "Expects error")
			s.Assert().ErrorIs(err, tc.expError, "Different error return")
		}
	}
}

func (s *SecretTestSuite) TestCreateInvalidSecretFolder() {
	type nameTestCase struct {
		name     string
		expOk    bool
		expError error
	}

	origContent := "original content"
	newContent := "new content"

	testcases := []nameTestCase{
		{name: "---", expOk: true, expError: nil},
		{name: "&#", expOk: false, expError: ErrBadPathName},
		{name: "", expOk: false, expError: ErrExists},
		{name: "   ", expOk: false, expError: ErrBadPathName},
		{name: "  /empty_parent", expOk: false, expError: ErrBadPathName},
		{name: "abc<def", expOk: false, expError: ErrBadPathName},

		// The following test cases are skipped now due to DV-231
		/*
			{name: "<", expOk: true, expError: nil},
			{name: ">", expOk: true, expError: nil},
			{name: "def>abc", expOk: true, expError: nil},
			{name: "&", expOk: true, expError: nil},
		*/
	}

	var secretName string
	objName := "test_file"
	for _, tc := range testcases {
		secretName = s.testFolder + "/" + tc.name
		s.T().Logf("For path [%s]:\n", secretName)
		success, _, _, err := s.handle.CreateFolder(secretName, "")
		if success {
			defer s.handle.Delete(secretName)
		}
		if tc.expOk {
			// expects success
			s.Assert().True(success, "Expects success")
			s.Assert().NoError(err, "Expects no error")

			// try to see if file operation works
			objPath := secretName + "/" + objName
			id := s.createTextSecret(objPath, origContent)
			defer s.handle.Delete(objPath)
			s.testModifySingleTextSecret(objPath, origContent, newContent, id)

			// try to list content of folder
			items, r, err := s.handle.List(secretName)
			s.Assert().NoError(err, "Should not get error when listing folder")
			s.Assert().Equal(200, r.StatusCode, "HTTP status should be 200")
			s.Assert().Equal(1, len(items), "Should return 1 item")
			expItem := Item{Name: objName, Type: SecretTypeText, ID: id}
			s.Assert().Equal(expItem, items[0], "Directory content")

		} else {
			// expects error
			s.Assert().False(success, "Expects error")
			s.Assert().ErrorIs(err, tc.expError, "Different error return")
		}
	}
}

// isPASenantConnected checks if a tenant is online and connected
// Note: use a new HTTP client but set the timeout to be 1 second
func isPASTenantConnected(tenantURL string) bool {

	cl := http.Client{
		Timeout: 1 * time.Second,
	}
	resp, err := cl.Get("https://" + tenantURL + "/health/ping")
	if err != nil {
		return false
	}
	// check response
	return resp.StatusCode == 200
}

// getAccessToken returns the Oauth access token for the user
func (s *SecretTestSuite) getAccessToken() (string, error) {
	// get oauth token for user
	oauthClient, err := oauth.GetNewConfidentialClient("https://"+s.Config.TenantURL, s.Config.PASuser.Username, s.Config.PASuser.Password, nil)
	if err != nil {
		s.T().Logf("Error in getting Oauth client: %v", err)
		return "", err
	}

	oauthToken, oauthError, err := oauthClient.ClientCredentials(s.Config.AppID, s.Config.Scope)
	if err != nil {
		s.T().Logf("Error in sending authentication request to server: %v", err)
		return "", err
	}

	if oauthError != nil {
		s.T().Logf("Authentication error: %v.  Description: %v\n", oauthError.Error, oauthError.Description)
		return "", errors.New(oauthError.Error)
	}
	return oauthToken.AccessToken, nil
}

// createTestFolder creates a test folder
// the caller expects the folder is created successfully
func (s *SecretTestSuite) createTestFolder(path string) string {

	success, id, r, err := s.handle.CreateFolder(path, "")
	s.Assert().NoErrorf(err, "CreateFolder [%s] should not result in error", path)
	s.Assert().NotEmptyf(id, "ID of created folder [%s] should be returned", path)
	s.Assert().Truef(success, "CreateFolder [%s] should return true for success", path)
	s.Assert().Equalf(201, r.StatusCode, "CreateFolder [%s] should return 201 status", path)
	return id
}

// createTextSecret creates a text secret
// the caller expects the secret is created successfully
func (s *SecretTestSuite) createTextSecret(path string, content string) string {
	success, id, r, err := s.handle.Create(path, "", content)
	s.Assert().NoErrorf(err, "Secret [%s] creation should return no error", path)
	s.Assert().NotEmptyf(id, "ID of secret [%s] must be returned", path)
	s.Assert().Truef(success, "Must return true for success when creating secret [%s]", path)
	s.Assert().Equalf(201, r.StatusCode, "CreateSecret [%s] should return HTTP status 201", path)
	return id // return on success
}

// createKVSecret creates a test keyvalue secret
// the caller expects the secret is created successfully
func (s *SecretTestSuite) createKVSecret(path string, content map[string]string) string {
	success, id, r, err := s.handle.Create(path, "", content)
	s.Assert().NoErrorf(err, "Secret [%s] creation should return no error", path)
	s.Assert().NotEmptyf(id, "ID of secret [%s] must be returned", path)
	s.Assert().Truef(success, "Must return true for success when creating secret [%s]", path)
	s.Assert().Equalf(201, r.StatusCode, "CreateSecret [%s] should return HTTP status 201", path)
	return id // return on success
}
