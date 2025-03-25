//go:build unit
// +build unit

package crypto

type CryptoManagerMock struct {
	err error
}

func NewCryptoManagerMock() *CryptoManagerMock {
	return &CryptoManagerMock{
		err: nil,
	}
}

func NewInvalidCryptoManagerMock(err error) *CryptoManagerMock {
	return &CryptoManagerMock{
		err: err,
	}
}

func (c *CryptoManagerMock) Encrypt(content []byte) (string, error) {
	return string(content), c.err
}

func (c *CryptoManagerMock) Decrypt(text string) ([]byte, error) {
	return []byte(text), c.err
}
