package auth

type Keys struct {
	privateKey         string
	publicKey          string
	secondaryPublicKey string
}

type KeysI interface {
}
