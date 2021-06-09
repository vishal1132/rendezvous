package rendezvous

// HashFunction is a type for hash functions for code readability.
// This is a key for scoreFuncs map
type HashName string

const (
	MD5    HashName = "md5"
	SHA256 HashName = "sha256"
	SHA1   HashName = "sha1"
)
