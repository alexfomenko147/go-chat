package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"golang.org/x/crypto/chacha20poly1305"
	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/hkdf"
)

type IdentityKeypair struct {
	PrivateKey ed25519.PrivateKey
	PublicKey  ed25519.PublicKey
}

func GenerateIdentityKeypair() (*IdentityKeypair, error) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate ed25519 keypair: %w", err)
	}
	return &IdentityKeypair{PrivateKey: priv, PublicKey: pub}, nil
}

func Fingerprint(pub ed25519.PublicKey) string {
	h := sha256.Sum256(pub)
	return hex.EncodeToString(h[:16])
}

type SessionKeys struct {
	SharedSecret []byte
	EncKey       []byte
	AuthKey      []byte
}

func DeriveSessionKeys(privateKey ed25519.PrivateKey, publicKey ed25519.PublicKey) (*SessionKeys, error) {
	curvePriv := ed25519PrivateKeyToCurve25519(privateKey)
	curvePub, err := ed25519PublicKeyToCurve25519(publicKey)
	if err != nil {
		return nil, fmt.Errorf("convert public key: %w", err)
	}

	shared, err := curve25519.X25519(curvePriv, curvePub)
	if err != nil {
		return nil, fmt.Errorf("x25519 shared secret: %w", err)
	}

	h := hkdf.New(sha256.New, shared, nil, []byte("go-chat-session-key"))
	encKey := make([]byte, 32)
	authKey := make([]byte, 32)
	if _, err := io.ReadFull(h, encKey); err != nil {
		return nil, fmt.Errorf("hkdf enc key: %w", err)
	}
	if _, err := io.ReadFull(h, authKey); err != nil {
		return nil, fmt.Errorf("hkdf auth key: %w", err)
	}

	return &SessionKeys{
		SharedSecret: shared,
		EncKey:       encKey,
		AuthKey:      authKey,
	}, nil
}

type Cipher struct {
	key []byte
}

func NewCipher(key []byte) *Cipher {
	return &Cipher{key: key}
}

func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(c.key)
	if err != nil {
		return nil, fmt.Errorf("new chacha20poly1305: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}

	ciphertext := aead.Seal(nil, nonce, plaintext, nil)
	return append(nonce, ciphertext...), nil
}

func (c *Cipher) Decrypt(data []byte) ([]byte, error) {
	aead, err := chacha20poly1305.New(c.key)
	if err != nil {
		return nil, fmt.Errorf("new chacha20poly1305: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}

	return plaintext, nil
}

func ed25519PrivateKeyToCurve25519(priv ed25519.PrivateKey) []byte {
	h := sha256.New()
	h.Write(priv[:32])
	digest := h.Sum(nil)
	digest[0] &= 248
	digest[31] &= 127
	digest[31] |= 64
	return digest
}

func ed25519PublicKeyToCurve25519(pub ed25519.PublicKey) ([]byte, error) {
	if len(pub) != ed25519.PublicKeySize {
		return nil, fmt.Errorf("invalid ed25519 public key length")
	}
	cp := make([]byte, 32)
	copy(cp, pub)
	return cp, nil
}
