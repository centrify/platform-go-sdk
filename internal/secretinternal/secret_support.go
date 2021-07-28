package secretinternal

// SecretMarshaler defines the interface that all SecretXXXWritable must support
// objects must support
type SecretMarshaler interface {
	MarshalJSON() ([]byte, error)
	UnmarshalJSON(bytes []byte) (err error)
}
