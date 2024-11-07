package auth

type Keys struct {
	PrivateKey         string
	PublicKey          string
	SecondaryPublicKey string
}

type KeysI interface {
}

var CreateKeys = func(privateKey string, publicKey string, secondaryPublicKey string) (KeysI, error) {
	return &Keys{
		privateKey,
		publicKey,
		secondaryPublicKey,
	}, nil
}
